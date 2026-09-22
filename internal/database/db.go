// Package database provides functionality for connecting to a PostgreSQL database using GORM.
package database

import (
	"fmt"

	"github.com/zekielmp/Bitly/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New establishes a new connection to the PostgreSQL database using the provided configuration settings.
// It returns a pointer to the gorm.DB instance and an error if the connection fails.
func New(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// func conn(name string) (*sql.DB, error) {

// 	db, err := sql.Open("pgx", name)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer db.Close()

// 	_ = db.Ping()
// 	return db, nil

// }
