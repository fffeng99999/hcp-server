package database

import (
	"fmt"
	"strings"

	"github.com/fffeng99999/hcp-server/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// 如果目标数据库不存在，则尝试自动创建数据库
		if strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "3D000") {
			if createErr := createDatabase(cfg); createErr != nil {
				return nil, fmt.Errorf("failed to create database: %w (original error: %v)", createErr, err)
			}
			// 创建成功后重新尝试连接
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return db, nil
}

func createDatabase(cfg config.DatabaseConfig) error {
	// 连接到默认的 postgres 数据库，用于创建业务数据库
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	// 检查目标数据库是否已经存在
	var exists bool
	checkStmt := fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s')", cfg.DBName)
	if err := db.Raw(checkStmt).Scan(&exists).Error; err != nil {
		return err
	}

	if !exists {
		// 创建目标业务数据库
		createStmt := fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)
		if err := db.Exec(createStmt).Error; err != nil {
			return err
		}
	}

	return nil
}
