package postgres

import (
	"fmt"

	"github.com/devusSs/crosshairs/config"
	"github.com/devusSs/crosshairs/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	tableUsers      = "user_accounts"
	tableCrosshairs = "crosshairs"
)

type psql struct {
	db *gorm.DB
}

func NewConnection(cfg *config.Config, gormLogger logger.Interface) (database.Service, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPassword,
		cfg.PostgresDB, cfg.PostgresPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, err
	}

	return &psql{db}, nil
}

func (p *psql) TestConnection() error {
	db, err := p.db.DB()
	if err != nil {
		return err
	}
	return db.Ping()
}

func (p *psql) CloseConnection() error {
	db, err := p.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (p *psql) MakeMigrations() error {
	if err := p.db.AutoMigrate(&database.UserAccount{}); err != nil {
		return err
	}
	return p.db.AutoMigrate(&database.Crosshair{})
}
