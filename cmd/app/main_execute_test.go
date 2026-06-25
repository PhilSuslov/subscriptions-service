package main

import (
	"context"
	"errors"
	"testing"

	"github.com/example/subscriptions-service/internal/config"
	"github.com/example/subscriptions-service/internal/repository/postgres"
)

func TestExecuteConfigError(t *testing.T) {
	oldLoad := configLoad
	t.Cleanup(func() { configLoad = oldLoad })
	configLoad = func(string) (*config.Config, error) { return nil, errors.New("boom") }
	if err := execute(); err == nil {
		t.Fatal("expected error")
	}
}

func TestRunPoolError(t *testing.T) {
	oldPool := newPool
	t.Cleanup(func() { newPool = oldPool })
	newPool = func(context.Context, config.PostgresConfig) (postgres.Queryer, error) {
		return nil, errors.New("boom")
	}
	cfg := &config.Config{}
	if err := run(cfg, nil); err == nil {
		t.Fatal("expected error")
	}
}
