package config

import (
	"sync"

	env "mall/app/Support/config"
)

type DatabaseConfig struct {
	Host         string
	Port         string
	Database     string
	Username     string
	Password     string
	Charset      string
	MaxOpenConns int
	MaxIdleConns int
}

var (
	dbCfg     *DatabaseConfig
	dbCfgOnce sync.Once
)

func Database() *DatabaseConfig {
	dbCfgOnce.Do(func() {
		dbCfg = &DatabaseConfig{
			Host:         env.Get("DB_HOST", "127.0.0.1"),
			Port:         env.Get("DB_PORT", "3306"),
			Database:     env.Get("DB_DATABASE", "my_mall"),
			Username:     env.Get("DB_USERNAME", "root"),
			Password:     env.Get("DB_PASSWORD", ""),
			Charset:      env.Get("DB_CHARSET", "utf8mb4"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 100),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 10),
		}
	})
	return dbCfg
}
