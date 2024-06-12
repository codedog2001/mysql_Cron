package Executor

import (
	"MySQL_Job/domain"
	"context"
)

type Executor interface {
	Name() string
	// Exec ctx 这个是全局控制，Executor 的实现者注意要正确处理 ctx 超时或者取消
	Exec(ctx context.Context, j domain.Job) error
}
