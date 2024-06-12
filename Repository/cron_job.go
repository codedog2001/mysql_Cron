package Repository

import (
	"MySQL_Job/Repository/DAO"
	"MySQL_Job/domain"
	"context"
	"time"
)

type CronJobRepository interface {
	Preempt(ctx context.Context) (domain.Job, error)
	Release(ctx context.Context, jid int64) error
	UpdateUtime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, time time.Time) error
	RecordJobExecution(ctx context.Context, jobId int64, start, end time.Time, status, message string) error
}

type PreemptJobRepository struct {
	dao DAO.JobDAO
}

func NewPreemptJobRepository(dao DAO.JobDAO) CronJobRepository {
	return &PreemptJobRepository{dao: dao}
}

func (p *PreemptJobRepository) Preempt(ctx context.Context) (domain.Job, error) {
	j, err := p.dao.Preempt(ctx)
	return domain.Job{
		Id:         j.Id,
		Name:       j.Name,
		Executor:   j.Executor,
		Expression: j.Expression,
		Cfg:        j.Cfg,
		Status:     j.Status,
		Version:    j.Version,
		NextTime:   j.NextTime,
		Utime:      j.Utime,
		Ctime:      j.Ctime,
	}, err
}

func (p *PreemptJobRepository) Release(ctx context.Context, jid int64) error {
	return p.dao.Release(ctx, jid)
}

func (p *PreemptJobRepository) UpdateUtime(ctx context.Context, id int64) error {
	return p.dao.UpdateUtime(ctx, id)
}

func (p *PreemptJobRepository) UpdateNextTime(ctx context.Context, id int64, time time.Time) error {
	return p.dao.UpdateNextTime(ctx, id, time)
}

func (p *PreemptJobRepository) RecordJobExecution(ctx context.Context, jobId int64, start, end time.Time, status, message string) error {
	return p.dao.RecordJobExecution(ctx, jobId, start, end, status, message)
}
