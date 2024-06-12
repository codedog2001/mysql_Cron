package Executor

import (
	"MySQL_Job/domain"
	"context"
)

type GrpcExecutor struct{}

func NewGrpcExecutor() *GrpcExecutor {
	return &GrpcExecutor{}
}

func (g *GrpcExecutor) Name() string {
	return "grpc"
}

func (g *GrpcExecutor) Exec(ctx context.Context, j domain.Job) error {
	// Implement gRPC call logic here
	return nil
}
