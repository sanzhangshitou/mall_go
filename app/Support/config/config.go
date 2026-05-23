package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Load 加载 .env 文件并返回配置 map
func Load(path string) error {
	return godotenv.Load(path)
}

// Get 获取字符串环境变量, 不存在时返回默认值
func Get(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

// GetInt 获取整数环境变量
func GetInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}

// GetBool 获取布尔环境变量
func GetBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	return b
}

// ---------- 预定义便捷方法 ----------

func AppName() string { return Get("APP_NAME", "Mall") }
func AppEnv() string  { return Get("APP_ENV", "local") }
func AppDebug() bool  { return GetBool("APP_DEBUG", true) }
func AppPort() string { return Get("APP_PORT", "8080") }

func DbHost() string      { return Get("DB_HOST", "127.0.0.1") }
func DbPort() string      { return Get("DB_PORT", "3306") }
func DbDatabase() string  { return Get("DB_DATABASE", "my_mall") }
func DbUsername() string  { return Get("DB_USERNAME", "root") }
func DbPassword() string  { return Get("DB_PASSWORD", "") }
func DbCharset() string   { return Get("DB_CHARSET", "utf8mb4") }
func DbMaxOpenConns() int { return GetInt("DB_MAX_OPEN_CONNS", 100) }
func DbMaxIdleConns() int { return GetInt("DB_MAX_IDLE_CONNS", 10) }

func RedisHost() string     { return Get("REDIS_HOST", "127.0.0.1") }
func RedisPort() string     { return Get("REDIS_PORT", "6379") }
func RedisPassword() string { return Get("REDIS_PASSWORD", "") }
func RedisDB() int          { return GetInt("REDIS_DB", 0) }
func RedisPoolSize() int    { return GetInt("REDIS_POOL_SIZE", 100) }

func LogPath() string    { return Get("LOG_PATH", "storage/logs/") }
func LogLevel() string   { return Get("LOG_LEVEL", "debug") }
func LogMaxSize() int    { return GetInt("LOG_MAX_SIZE", 100) }
func LogMaxBackups() int { return GetInt("LOG_MAX_BACKUPS", 10) }
func LogMaxAge() int     { return GetInt("LOG_MAX_AGE", 30) }

func PageDefaultSize() int { return GetInt("PAGE_DEFAULT_SIZE", 15) }
func PageMaxSize() int     { return GetInt("PAGE_MAX_SIZE", 100) }
