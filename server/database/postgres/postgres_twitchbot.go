package postgres

import "github.com/devusSs/crosshairs/database"

func (p *psql) WriteTwitchBotLog(botLog *database.TwitchBotLog) error {
	tx := p.db.Table("twitch_bot_logs").Create(&botLog)
	return tx.Error
}

func (p *psql) GetAllTwitchBotLogEntries() ([]*database.TwitchBotLog, error) {
	var logs []*database.TwitchBotLog
	tx := p.db.Table("twitch_bot_logs").Find(&logs)
	return logs, tx.Error
}

func (p *psql) GetLatestTwitchBotLogWithLimit(limit int) ([]*database.TwitchBotLog, error) {
	var logs []*database.TwitchBotLog
	tx := p.db.Table("twitch_bot_logs").Order("created_at desc").Limit(limit).Find(&logs)
	return logs, tx.Error
}

func (p *psql) GetLatestTwitchBotLogByType(logType string) ([]*database.TwitchBotLog, error) {
	var logs []*database.TwitchBotLog
	tx := p.db.Table("twitch_bot_logs").Order("created_at desc").Where("issuer = ?", logType).Find(&logs)
	return logs, tx.Error
}

func (p *psql) GetLatestTwitchBotLogByTypeWithLimit(logType string, limit int) ([]*database.TwitchBotLog, error) {
	var logs []*database.TwitchBotLog
	tx := p.db.Table("twitch_bot_logs").Order("created_at desc").Where("issuer = ?", logType).Limit(limit).Find(&logs)
	return logs, tx.Error
}

func (p *psql) AddTwitchTokenRefreshStore(token *database.TwitchRefreshTokenStore) (*database.TwitchRefreshTokenStore, error) {
	tx := p.db.Table("twitch_refresh_token_stores").Create(&token)
	return token, tx.Error
}

func (p *psql) GetLatestTwitchTokenRefreshStore(loginName *database.TwitchRefreshTokenStore) (*database.TwitchRefreshTokenStore, error) {
	tx := p.db.Table("twitch_refresh_token_stores").Order("created_at desc").Where("twitch_login = ?", loginName.TwitchLogin).First(&loginName)
	return loginName, tx.Error
}

func (p *psql) DeleteAllTwitchTokenRefreshStore(loginName *database.TwitchRefreshTokenStore) error {
	tx := p.db.Table("twitch_refresh_token_stores").Where("twitch_login = ?", loginName.TwitchLogin).Delete(&database.TwitchRefreshTokenStore{})
	return tx.Error
}
