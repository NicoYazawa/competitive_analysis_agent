package storage

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupMiniRedis(t *testing.T) (*RedisCache, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to create miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	cache := &RedisCache{
		client: client,
	}

	return cache, mr
}

func TestRedisCache_SetAndGet(t *testing.T) {
	cache, mr := setupMiniRedis(t)
	defer mr.Close()

	ctx := context.Background()

	err := cache.Set(ctx, "key1", "value1", time.Hour)
	assert.NoError(t, err)

	val, err := cache.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)
}

func TestRedisCache_Get_NotFound(t *testing.T) {
	cache, mr := setupMiniRedis(t)
	defer mr.Close()

	ctx := context.Background()

	_, err := cache.Get(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
}

func TestRedisCache_Delete(t *testing.T) {
	cache, mr := setupMiniRedis(t)
	defer mr.Close()

	ctx := context.Background()

	cache.Set(ctx, "key1", "value1", time.Hour)
	err := cache.Delete(ctx, "key1")
	assert.NoError(t, err)

	_, err = cache.Get(ctx, "key1")
	assert.Equal(t, redis.Nil, err)
}

func TestRedisCache_Exists(t *testing.T) {
	cache, mr := setupMiniRedis(t)
	defer mr.Close()

	ctx := context.Background()

	cache.Set(ctx, "key1", "value1", time.Hour)

	exists, err := cache.Exists(ctx, "key1")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = cache.Exists(ctx, "key2")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestRedisCache_TTL(t *testing.T) {
	cache, mr := setupMiniRedis(t)
	defer mr.Close()

	ctx := context.Background()

	cache.Set(ctx, "key1", "value1", time.Hour)

	ttl, err := cache.TTL(ctx, "key1")
	assert.NoError(t, err)
	assert.True(t, ttl > 0 && ttl <= time.Hour)
}

func TestRedisCache_Expire(t *testing.T) {
	cache, mr := setupMiniRedis(t)
	defer mr.Close()

	ctx := context.Background()

	cache.Set(ctx, "key1", "value1", time.Hour)
	err := cache.Expire(ctx, "key1", time.Minute)
	assert.NoError(t, err)

	ttl, err := cache.TTL(ctx, "key1")
	assert.NoError(t, err)
	assert.True(t, ttl > 0 && ttl <= time.Minute)
}

func TestRedisCache_Incr(t *testing.T) {
	cache, mr := setupMiniRedis(t)
	defer mr.Close()

	ctx := context.Background()

	val, err := cache.Incr(ctx, "counter")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), val)

	val, err = cache.Incr(ctx, "counter")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), val)
}

func TestRedisCache_IncrBy(t *testing.T) {
	cache, mr := setupMiniRedis(t)
	defer mr.Close()

	ctx := context.Background()

	val, err := cache.IncrBy(ctx, "counter", 5)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), val)

	val, err = cache.IncrBy(ctx, "counter", 3)
	assert.NoError(t, err)
	assert.Equal(t, int64(8), val)
}

func TestRedisCache_ZAdd(t *testing.T) {
	cache, mr := setupMiniRedis(t)
	defer mr.Close()

	ctx := context.Background()

	members := []redis.Z{
		{Score: 1, Member: "a"},
		{Score: 2, Member: "b"},
		{Score: 3, Member: "c"},
	}

	err := cache.ZAdd(ctx, "sortedset", members...)
	assert.NoError(t, err)

	result, err := cache.ZRangeByScore(ctx, "sortedset", &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	})
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, result)
}

func TestRedisCache_ZRangeByScore(t *testing.T) {
	cache, mr := setupMiniRedis(t)
	defer mr.Close()

	ctx := context.Background()

	members := []redis.Z{
		{Score: 1, Member: "a"},
		{Score: 2, Member: "b"},
		{Score: 3, Member: "c"},
	}
	cache.ZAdd(ctx, "sortedset", members...)

	result, err := cache.ZRangeByScore(ctx, "sortedset", &redis.ZRangeBy{
		Min: "2",
		Max: "+inf",
	})
	assert.NoError(t, err)
	assert.Equal(t, []string{"b", "c"}, result)
}

func TestRedisCache_Close(t *testing.T) {
	cache, mr := setupMiniRedis(t)
	mr.Close()

	err := cache.Close()
	assert.NoError(t, err)
}

func TestRedisCache_HealthCheck(t *testing.T) {
	cache, mr := setupMiniRedis(t)
	defer mr.Close()

	ctx := context.Background()
	err := cache.HealthCheck(ctx)
	assert.NoError(t, err)
}
