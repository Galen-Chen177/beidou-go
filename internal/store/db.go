package store

import (
	"fmt"

	"beidou-go/config"
	"beidou-go/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// DB 返回全局数据库实例
func DB() *gorm.DB {
	return db
}

// InitDB 初始化数据库连接并自动迁移
func InitDB(cfg config.DatabaseConfig) error {
	var err error
	db, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // 生产模式只记录警告
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层 sql.DB 失败: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	return nil
}

// AutoMigrate 自动迁移表结构
func AutoMigrate() error {
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return db.AutoMigrate(
		&model.Account{},
		&model.Character{},
		&model.Item{},
		&model.Skill{},
	)
}
