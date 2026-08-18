package storage

import (
	"context"
	"fmt"
	"time"

	"competitive-analysis-agent/internal/config"

	"github.com/redis/go-redis/v9"
)

// RedisCache Redis缓存管理
type RedisCache struct {
	client *redis.Client
	cfg    *config.RedisConfig
}

// NewRedisCache 创建Redis缓存实例
func NewRedisCache(cfg *config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisCache{client: client, cfg: cfg}, nil
}

// Client 获取redis客户端
func (r *RedisCache) Client() *redis.Client {
	return r.client
}

// Set 设置键值对
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// SetBytes 设置字节数据
func (r *RedisCache) SetBytes(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Get 获取值
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// GetBytes 获取字节数据
func (r *RedisCache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

// Delete 删除键
func (r *RedisCache) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	return n > 0, err
}

// TTL 获取键的TTL
func (r *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// Expire 设置键的过期时间
func (r *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.client.Expire(ctx, key, ttl).Err()
}

// Incr 自增
func (r *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// IncrBy 增加指定值
func (r *RedisCache) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.IncrBy(ctx, key, value).Result()
}

// ZAdd 添加有序集合成员
func (r *RedisCache) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return r.client.ZAdd(ctx, key, members...).Err()
}

// ZRangeByScore 按分数范围获取有序集合
func (r *RedisCache) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return r.client.ZRangeByScore(ctx, key, opt).Result()
}

// ZRevRangeByScore 逆序按分数范围获取有序集合
func (r *RedisCache) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return r.client.ZRevRangeByScore(ctx, key, opt).Result()
}

// Close 关闭连接
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// HealthCheck 健康检查
func (r *RedisCache) HealthCheck(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
