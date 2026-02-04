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
		// If database does not exist, try to create it
		if strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "3D000") {
			if createErr := createDatabase(cfg); createErr != nil {
				return nil, fmt.Errorf("failed to create database: %w (original error: %v)", createErr, err)
			}
			// Retry connection
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
	// Connect to 'postgres' database to create the new database
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

	// Check if database exists
	var exists bool
	checkStmt := fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s')", cfg.DBName)
	if err := db.Raw(checkStmt).Scan(&exists).Error; err != nil {
		return err
	}

	if !exists {
		// Create database
		createStmt := fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)
		if err := db.Exec(createStmt).Error; err != nil {
			return err
		}
	}

	return nil
}
