package config

import (
	"sync"

	env "mall/app/Support/config"
)

type LogConfig struct {
	Path       string
	Level      string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

var (
	logCfg     *LogConfig
	logCfgOnce sync.Once
)

func Log() *LogConfig {
	logCfgOnce.Do(func() {
		logCfg = &LogConfig{
			Path:       env.Get("LOG_PATH", "storage/logs/"),
			Level:      env.Get("LOG_LEVEL", "debug"),
			MaxSize:    env.GetInt("LOG_MAX_SIZE", 100),
			MaxBackups: env.GetInt("LOG_MAX_BACKUPS", 10),
			MaxAge:     env.GetInt("LOG_MAX_AGE", 30),
		}
	})
	return logCfg
}
