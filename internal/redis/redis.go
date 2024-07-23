package redis

import (
	"context"
	"fmt"
	"grpc-redis-postgres/proto"

	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	Get(ctx context.Context, key string) (map[string]string, error)
	Set(ctx context.Context, key string, user *proto.User) error
}

type redisClient struct {
	client *redis.Client
}

func NewRedisClient(addr string, password string, db int) (RedisClient, error) {
	fmt.Println("addr:", addr)
	fmt.Println("password:", password)
	fmt.Println("db:", db)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("error connecting to Redis: %w", err)
	}
	return &redisClient{client: client}, nil
}

func (r *redisClient) Get(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}
func (r *redisClient) Set(ctx context.Context, key string, user *proto.User) error {
	data := make(map[string]interface{})
	data["id"] = user.Id
	data["name"] = user.Name
	data["email"] = user.Email
	return r.client.HMSet(ctx, key, data).Err()
}
