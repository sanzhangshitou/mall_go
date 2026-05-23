package config

import (
	"sync"

	env "mall/app/Support/config"
)

type PageConfig struct {
	DefaultSize int
	MaxSize     int
}

var (
	pageCfg     *PageConfig
	pageCfgOnce sync.Once
)

func Page() *PageConfig {
	pageCfgOnce.Do(func() {
		pageCfg = &PageConfig{
			DefaultSize: env.GetInt("PAGE_DEFAULT_SIZE", 15),
			MaxSize:     env.GetInt("PAGE_MAX_SIZE", 100),
		}
	})
	return pageCfg
}
