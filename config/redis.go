package config

import (
	"sync"

	env "mall/app/Support/config"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

var (
	redisCfg     *RedisConfig
	redisCfgOnce sync.Once
)

func Redis() *RedisConfig {
	redisCfgOnce.Do(func() {
		redisCfg = &RedisConfig{
			Host:     env.Get("REDIS_HOST", "127.0.0.1"),
			Port:     env.Get("REDIS_PORT", "6379"),
			Password: env.Get("REDIS_PASSWORD", ""),
			DB:       env.GetInt("REDIS_DB", 0),
			PoolSize: env.GetInt("REDIS_POOL_SIZE", 100),
		}
	})
	return redisCfg
}
