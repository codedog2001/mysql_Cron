package Service

import (
	"MySQL_Job/Repository"
	"MySQL_Job/domain"
	"MySQL_Job/pkg/logger"
	"context"
	"time"
)

type CronJobService interface {
	Preempt(ctx context.Context) (domain.Job, error)
	ResetNextTime(ctx context.Context, j domain.Job) error
	RecordJobExecution(ctx context.Context, jobId int64, start, end time.Time, status string, err error) error
}

type cronJobService struct {
	repo Repository.CronJobRepository
	l    logger.LoggerV1
}

func NewCronJobService(repo Repository.CronJobRepository, l logger.LoggerV1) CronJobService {
	return &cronJobService{repo: repo, l: l}
}

func (c *cronJobService) Preempt(ctx context.Context) (domain.Job, error) {
	//首先进行抢占
	j, err := c.repo.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	//一分钟续约一次
	ticker := time.NewTicker(time.Minute)
	go func() {
		//利用tiker进行续约
		for range ticker.C {
			c.refresh(j.Id)
		}
	}()
	//取消函数，抢占任务之后就要定义好取消函数，方便后序schedler进行调用
	j.CancelFunc = func() {
		ticker.Stop() //停止计时，这个tiker是用来续约的，不是用来定时启动任务的
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		//释放job
		err := c.repo.Release(ctx, j.Id)
		if err != nil {
			c.l.Error("释放 job 失败",
				logger.Error(err),
				logger.Int64("jid", j.Id))
		}
	}
	return j, err
}

// ResetNextTime 更新下一次的启动时间
func (c *cronJobService) ResetNextTime(ctx context.Context, j domain.Job) error {
	//nexttimefunc在定义job的时候就要一起定义好
	nextTime := j.NextTimeFunc()
	return c.repo.UpdateNextTime(ctx, j.Id, nextTime)
}

// RecordJobExecution 任务执行的历史记录
func (c *cronJobService) RecordJobExecution(ctx context.Context, jobId int64, start, end time.Time, status string, err error) error {
	message := ""
	if err != nil {
		message = err.Error()
	}
	return c.repo.RecordJobExecution(ctx, jobId, start, end, status, message)
}

func (c *cronJobService) refresh(id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := c.repo.UpdateUtime(ctx, id)
	if err != nil {
		c.l.Error("续约失败", logger.Error(err),
			logger.Int64("jid", id))
	}
}
