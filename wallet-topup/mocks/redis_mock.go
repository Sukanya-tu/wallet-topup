package mocks

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

type RedisMock struct {
	mock.Mock
}

func (m *RedisMock) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)

	return redis.NewStringResult(args.String(0), args.Error(1))
}

func (m *RedisMock) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)

	return redis.NewStatusResult("", args.Error(0))
}

func (m *RedisMock) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)

	return redis.NewIntResult(1, args.Error(0))
}
