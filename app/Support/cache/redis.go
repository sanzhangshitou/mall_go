package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"mall/app/Support/logger"
	"mall/config"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	rdb  *redis.Client
	once sync.Once
	ctx  = context.Background()
)

// Init 初始化 Redis 连接
func Init() {
	once.Do(func() {
		rdb = newRedis()
	})
}

// Client 获取 Redis 客户端
func Client() *redis.Client {
	return rdb
}

// Ctx 获取默认上下文
func Ctx() context.Context {
	return ctx
}

func newRedis() *redis.Client {
	cfg := config.Redis()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Fatal("Redis 连接失败", zap.Error(err))
	}

	logger.Info("Redis 连接成功", zap.String("addr", client.Options().Addr))
	return client
}

// ---------- 便捷方法 ----------

// Set 设置缓存
func Set(key string, value interface{}, expiration time.Duration) error {
	return rdb.Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
func Get(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

// Del 删除缓存
func Del(keys ...string) error {
	return rdb.Del(ctx, keys...).Err()
}

// Exists 检查 key 是否存在
func Exists(key string) (bool, error) {
	n, err := rdb.Exists(ctx, key).Result()
	return n > 0, err
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) error {
	return rdb.Expire(ctx, key, expiration).Err()
}

// Incr 自增
func Incr(key string) (int64, error) {
	return rdb.Incr(ctx, key).Result()
}

// HSet 设置哈希字段
func HSet(key string, values ...interface{}) error {
	return rdb.HSet(ctx, key, values...).Err()
}

// HGet 获取哈希字段
func HGet(key, field string) (string, error) {
	return rdb.HGet(ctx, key, field).Result()
}

// HGetAll 获取哈希所有字段
func HGetAll(key string) (map[string]string, error) {
	return rdb.HGetAll(ctx, key).Result()
}

// CacheKey 生成统一前缀的缓存 key
func CacheKey(parts ...string) string {
	key := "mall"
	for _, p := range parts {
		key += ":" + p
	}
	return key
}
