package Executor

import (
	"MySQL_Job/domain"
	"context"
	"errors"
	"net/http"
)

type HttpExecutor struct{}

func NewHttpExecutor() *HttpExecutor {
	return &HttpExecutor{}
}

func (h *HttpExecutor) Name() string {
	return "http"
}

func (h *HttpExecutor) Exec(ctx context.Context, j domain.Job) error {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", j.Cfg, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("HTTP request failed with status " + resp.Status)
	}
	return nil
}
