package DAO

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type JobDAO interface {
	Preempt(ctx context.Context) (Job, error)
	Release(ctx context.Context, jid int64) error
	UpdateUtime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, t time.Time) error
	RecordJobExecution(ctx context.Context, jobId int64, start, end time.Time, status, message string) error
}

type GORMJobDAO struct {
	db *gorm.DB
}

func NewGORMJobDAO(db *gorm.DB) JobDAO {
	return &GORMJobDAO{db: db}
}

func (dao *GORMJobDAO) Preempt(ctx context.Context) (Job, error) {
	db := dao.db.WithContext(ctx)
	for {
		var j Job
		now := time.Now().UnixMilli()
		//先去数据库中找是等待的状态，并且下一次调度的时间已经到了
		err := db.Where("status = ? AND next_time <?",
			jobStatusWaiting, now).
			First(&j).Error
		if err != nil {
			return j, err
		}
		//拿到任务后，再去修改任务状态，修改成功后才能视为已经拿到了这个任务
		//version存在的作用是实现了乐观锁，即更新数据之前要先判断数据是否被别人修改了
		//高并发的场景下，如果有多个进程同时达到这里，那么先修改version的就会拿到锁，
		//而后修改version的，会因为jid和jversion联合找不到对应的数据，所以会导致更新失败，即rowaffected==0
		//这就是乐观锁，修改的时候再去判断，如果遇到同时修改的情况，那么就放弃修改
		res := db.WithContext(ctx).Model(&Job{}).
			Where("id = ? AND version = ?", j.Id, j.Version).
			Updates(map[string]any{
				"status":  jobStatusRunning,
				"version": j.Version + 1,
				"utime":   now,
			})
		if res.Error != nil {
			return Job{}, res.Error
		}
		if res.RowsAffected == 0 {
			continue
		} //如果version对不上，那么这里就是0
		return j, err
	}
}

// 记
func (dao *GORMJobDAO) Release(ctx context.Context, jid int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", jid).Updates(map[string]any{
		"status": jobStatusWaiting,
		"utime":  now,
	}).Error
}

func (dao *GORMJobDAO) UpdateUtime(ctx context.Context, jid int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", jid).Updates(map[string]any{
		"utime": now,
	}).Error
}

func (dao *GORMJobDAO) UpdateNextTime(ctx context.Context, jid int64, t time.Time) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", jid).Updates(map[string]any{
		"utime":     now,
		"next_time": t.UnixMilli(),
	}).Error
}

func (dao *GORMJobDAO) RecordJobExecution(ctx context.Context, jobId int64, start, end time.Time, status, message string) error {
	return dao.db.WithContext(ctx).Create(&JobExecutionHistory{
		JobId:     jobId,
		StartTime: start,
		EndTime:   end,
		Status:    status,
		Message:   message,
	}).Error
}
