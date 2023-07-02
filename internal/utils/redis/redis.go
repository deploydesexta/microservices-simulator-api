package redis

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	host                   = os.Getenv("REDIS_HOST")
	pass                   = os.Getenv("REDIS_PASS")
	Infinity time.Duration = 0
)

type (
	Pipe = redis.Pipeliner

	Service interface {
		Incr(ctx context.Context, key string) (int64, error)
		Expire(ctx context.Context, key string, duration time.Duration) (bool, error)
		Get(ctx context.Context, key string) (string, error)
		Set(ctx context.Context, key string, value interface{}) error
		HSet(ctx context.Context, key string, value interface{}) error
		HGet(ctx context.Context, key string, field string) (string, error)
		SAdd(ctx context.Context, key string, field string) error
		SMembers(ctx context.Context, key string) ([]string, error)
		TxPipelined(ctx context.Context, fn func(Pipe) error) ([]redis.Cmder, error)
		Pipelined(ctx context.Context, fn func(Pipe) error) ([]redis.Cmder, error)
	}

	Client struct {
		client *redis.Client
	}
)

func NewClient() *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pass, // no password set
		DB:       0,    // use default DB
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}

	return &Client{rdb}
}

func (r *Client) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *Client) Expire(ctx context.Context, key string, duration time.Duration) (bool, error) {
	return r.client.Expire(ctx, key, duration).Result()
}

func (r *Client) Get(ctx context.Context, key string) (string, error) {
	cmd := r.client.Get(ctx, key)
	if cmd.Err() == redis.Nil {
		return "", nil
	}
	return cmd.Result()
}

func (r *Client) Set(ctx context.Context, key string, value interface{}) error {
	return r.client.Set(ctx, key, value, Infinity).Err()
}

func (r *Client) HSet(ctx context.Context, key string, value interface{}) error {
	return r.client.HSet(ctx, key, value, Infinity).Err()
}

func (r *Client) HGet(ctx context.Context, key string, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

func (r *Client) SAdd(ctx context.Context, key string, field string) error {
	return r.client.SAdd(ctx, key, field).Err()
}

func (r *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

func (r *Client) TxPipelined(ctx context.Context, fn func(Pipe) error) ([]redis.Cmder, error) {
	return r.client.TxPipelined(ctx, fn)
}

func (r *Client) Pipelined(ctx context.Context, fn func(Pipe) error) ([]redis.Cmder, error) {
	return r.client.Pipelined(ctx, fn)
}
