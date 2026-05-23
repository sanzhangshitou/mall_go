package config

import (
	"sync"

	env "mall/app/Support/config"
)

type AppConfig struct {
	Name  string
	Env   string
	Debug bool
	Port  string
}

var (
	appCfg     *AppConfig
	appCfgOnce sync.Once
)

func App() *AppConfig {
	appCfgOnce.Do(func() {
		appCfg = &AppConfig{
			Name:  env.Get("APP_NAME", "Mall"),
			Env:   env.Get("APP_ENV", "local"),
			Debug: env.GetBool("APP_DEBUG", true),
			Port:  env.Get("APP_PORT", "8080"),
		}
	})
	return appCfg
}
