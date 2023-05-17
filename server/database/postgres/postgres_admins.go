package postgres

import "github.com/devusSs/crosshairs/database"

func (p *psql) GetAllUsers() ([]*database.UserAccount, error) {
	var users []*database.UserAccount
	tx := p.db.Table(tableUsers).Find(&users)
	return users, tx.Error
}

func (p *psql) GetAllCrosshairs() ([]*database.Crosshair, error) {
	var crosshairs []*database.Crosshair
	tx := p.db.Table(tableCrosshairs).Find(&crosshairs)
	return crosshairs, tx.Error
}
