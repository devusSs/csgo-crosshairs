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
	tableEvents     = "events"
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
		TranslateError:         true,
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
	if err := p.db.AutoMigrate(&database.Crosshair{}); err != nil {
		return err
	}
	if err := p.db.AutoMigrate(&database.Event{}); err != nil {
		return err
	}
	if err := p.db.AutoMigrate(&database.EngineerToken{}); err != nil {
		return err
	}
	if err := p.db.AutoMigrate(&database.TwitchBotLog{}); err != nil {
		return err
	}
	return p.db.AutoMigrate(&database.TwitchRefreshTokenStore{})
}

func (p *psql) GetPostgresVersion() (string, error) {
	db, err := p.db.DB()
	if err != nil {
		return "", err
	}

	var version string
	err = db.QueryRow("select version()").Scan(&version)

	return version, err
}

func (p *psql) CreateNewEngineerToken(token string) error {
	tx := p.db.Table("engineer_tokens").Create(&database.EngineerToken{Token: token})
	return tx.Error
}

func (p *psql) GetLatestEngineerToken() (string, error) {
	var token database.EngineerToken
	tx := p.db.Table("engineer_tokens").Last(&token)
	return token.Token, tx.Error
}
