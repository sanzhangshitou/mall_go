package bootstrap

import (
	"mall/app/Support/cache"
	env "mall/app/Support/config"
	"mall/app/Support/database"
	"mall/app/Support/logger"
	"mall/config"
)

// App 应用容器, 持有所有已初始化的服务
type App struct{}

// Boot 启动所有核心服务: 配置 -> 日志 -> 数据库 -> Redis
func Boot(envPath string) {
	// 1. 加载环境变量
	if err := env.Load(envPath); err != nil {
		panic("加载 .env 失败: " + err.Error())
	}

	// 2. 初始化日志
	logger.Init()

	// 3. 初始化数据库
	database.Init()

	// 4. 初始化 Redis
	cache.Init()

	logger.Info("应用启动完成 [" + config.App().Name + "]")
}
