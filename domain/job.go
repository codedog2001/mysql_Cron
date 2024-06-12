package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Job struct {
	Id   int64
	Name string
	// Cron 表达式
	Expression string
	Executor   string //这个任务规定了相应了执行器，grpc/http
	Cfg        string
	CancelFunc func()
	Status     int
	Version    int
	NextTime   int64
	Utime      int64
	Ctime      int64
}

func (j Job) NextTimeFunc() time.Time {
	c := cron.NewParser(cron.Second | cron.Minute | cron.Hour |
		cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	s, _ := c.Parse(j.Expression)
	return s.Next(time.Now())
}
