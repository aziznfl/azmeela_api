package repository

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/azmeela/sispeg-api/internal/domain"
	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	client *redis.Client
}

// NewRedisRepository will create an object that represents the domain.RedisRepository interface
func NewRedisRepository(client *redis.Client) domain.RedisRepository {
	return &redisRepository{client}
}

func (r *redisRepository) StoreRefreshToken(ctx context.Context, userID int, token string, duration time.Duration) error {
	return r.client.Set(ctx, token, userID, duration).Err()
}

func (r *redisRepository) GetRefreshToken(ctx context.Context, token string) (int, error) {
	val, err := r.client.Get(ctx, token).Result()
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *redisRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.client.Del(ctx, token).Err()
}

func (r *redisRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *redisRepository) Get(ctx context.Context, key string, target interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), target)
}

func (r *redisRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *redisRepository) FlushAll(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
}
