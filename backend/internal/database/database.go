package database

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/yoshsoftdev/porto/backend/internal/config"
	"github.com/yoshsoftdev/porto/backend/internal/models"
)

// Open connects to PostgreSQL when DATABASE_URL is set, otherwise to a local SQLite file.
func Open(cfg config.Config) (*gorm.DB, error) {
	gcfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)}

	var (
		db  *gorm.DB
		err error
	)
	if cfg.DatabaseURL != "" {
		log.Println("database: postgres")
		db, err = gorm.Open(postgres.Open(cfg.DatabaseURL), gcfg)
	} else {
		log.Printf("database: sqlite (%s)", cfg.SQLitePath)
		db, err = gorm.Open(sqlite.Open(cfg.SQLitePath+"?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)"), gcfg)
		if err == nil {
			if sqlDB, e := db.DB(); e == nil {
				sqlDB.SetMaxOpenConns(1)
			}
		}
	}
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(models.All()...); err != nil {
		return nil, err
	}
	return db, nil
}
