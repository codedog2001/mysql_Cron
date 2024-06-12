package Scheduler

import (
	"MySQL_Job/Service"
	"MySQL_Job/domain"
	"MySQL_Job/pkg/logger"
	"context"
	"golang.org/x/sync/semaphore"
	"time"
)

type Executor interface {
	Name() string
	// Exec ctx 这个是全局控制，Executor 的实现者注意要正确处理 ctx 超时或者取消
	Exec(ctx context.Context, j domain.Job) error
}

type Scheduler struct {
	dbTimeout time.Duration
	svc       Service.CronJobService
	permSvc   Service.PermissionService
	executors map[string]Executor
	l         logger.LoggerV1
	limiter   *semaphore.Weighted
}

func NewScheduler(svc Service.CronJobService, permSvc Service.PermissionService, l logger.LoggerV1) *Scheduler {
	return &Scheduler{
		svc:       svc,
		permSvc:   permSvc,
		dbTimeout: time.Second,
		limiter:   semaphore.NewWeighted(100),
		l:         l,
		executors: map[string]Executor{},
	}
}

// RegisterExecutor 通过这个函数来注册执行器，hhtp or grpc
// 需要先new一个scheduler再进行注册
func (s *Scheduler) RegisterExecutor(exec Executor) {
	s.executors[exec.Name()] = exec
}

// Schedule 鉴权，调度，历史记录，
func (s *Scheduler) Schedule(ctx context.Context, userID int64, requiredPermission string) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		//最多允许一百个节点调度
		err := s.limiter.Acquire(ctx, 1)
		if err != nil {
			return err
		}
		//鉴权
		hasPermission, err := s.permSvc.CheckPermission(ctx, userID, requiredPermission)
		if err != nil || !hasPermission {
			//没有权限或者鉴权失败，直接释放limiter，且返回，不用继续往下l
			s.limiter.Release(1)
			return err
		}
		//在规定时间内抢占任务
		dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
		j, err := s.svc.Preempt(dbCtx)
		cancel()
		if err != nil {
			continue
		} //如果没抢到，就继续抢
		exec, ok := s.executors[j.Executor]
		if !ok {
			s.l.Error("找不到执行器",
				logger.Int64("jid", j.Id),
				logger.String("executor", j.Executor))
			continue
		} //走到这说明拿到任务，且找到了相应的执行器
		go func() {
			defer func() {
				s.limiter.Release(1)
				j.CancelFunc()
				//这里defer语句可以保证job执行完了，才会被取消和释放
				//只有主动释放，别人才能拿到这个任务
			}()
			start := time.Now()
			//执行，并记录起始时间
			err1 := exec.Exec(ctx, j)
			end := time.Now()
			status := "success"
			if err1 != nil {
				s.l.Error("执行任务失败",
					logger.Int64("jid", j.Id),
					logger.Error(err1))
				status = "failure"
			}
			err2 := s.svc.RecordJobExecution(ctx, j.Id, start, end, status, err1)
			if err2 != nil {
				s.l.Error("记录任务执行历史失败",
					logger.Int64("jid", j.Id),
					logger.Error(err2))
			} //执行完任务且记录历史记录后，就计算该任务的下一次调度
			err1 = s.svc.ResetNextTime(ctx, j)
			if err1 != nil {
				s.l.Error("重置下次执行时间失败",
					logger.Int64("jid", j.Id),
					logger.Error(err1))
			}
		}()
	}
}
