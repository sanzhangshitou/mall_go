package database

import (
	"fmt"
	"sync"
	"time"

	"mall/app/Support/logger"
	"mall/config"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

// Init 初始化 MySQL 连接
func Init() {
	once.Do(func() {
		db = newDB()
	})
}

// DB 获取数据库实例
func DB() *gorm.DB {
	return db
}

func newDB() *gorm.DB {
	cfg := config.Database()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
	)

	// GORM 日志级别
	logLevel := gormLogger.Silent
	if config.App().Debug {
		logLevel = gormLogger.Info
	}

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(logLevel),
	})
	if err != nil {
		logger.Fatal("数据库连接失败", zap.Error(err))
	}

	sqlDB, err := database.DB()
	if err != nil {
		logger.Fatal("获取 sql.DB 失败", zap.Error(err))
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("数据库连接成功", zap.String("host", cfg.Host))

	return database
}
