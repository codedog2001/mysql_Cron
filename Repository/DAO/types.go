package DAO

import "time"

type Job struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	Name       string `gorm:"type:varchar(128);unique"`
	Executor   string
	Expression string
	Cfg        string
	Status     int
	Version    int
	NextTime   int64 `gorm:"index"`
	Utime      int64
	Ctime      int64
}

type JobExecutionHistory struct {
	Id        int64 `gorm:"primaryKey,autoIncrement"`
	JobId     int64 `gorm:"index"`
	StartTime time.Time
	EndTime   time.Time
	Status    string
	Message   string
}
type Department struct {
	Id          int64  `gorm:"primaryKey,autoIncrement"`
	Name        string `gorm:"type:varchar(128);not null"`
	Description string `gorm:"type:text"`
}
type User struct {
	Id           int64      `gorm:"primaryKey,autoIncrement"`
	Username     string     `gorm:"type:varchar(128);not null"`
	Password     string     `gorm:"type:varchar(256);not null"`
	DepartmentId int64      `gorm:"index"`
	Department   Department `gorm:"foreignKey:DepartmentId"`
}
type Permission struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	UserId     int64  `gorm:"index"`
	Permission string `gorm:"type:varchar(128)"`
	User       User   `gorm:"foreignKey:UserId"`
}

const (
	jobStatusWaiting = iota
	jobStatusRunning
	jobStatusPaused
)
