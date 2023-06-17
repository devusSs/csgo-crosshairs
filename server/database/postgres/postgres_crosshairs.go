package postgres

import (
	"github.com/devusSs/crosshairs/database"
	"github.com/google/uuid"
)

func (p *psql) AddCrosshair(ch *database.Crosshair) (*database.Crosshair, error) {
	tx := p.db.Table(tableCrosshairs).Save(ch)
	return ch, tx.Error
}

func (p *psql) GetAllCrosshairsFromUser(user uuid.UUID) ([]*database.Crosshair, error) {
	var crosshairs []*database.Crosshair
	tx := p.db.Table(tableCrosshairs).Where("registrant_id = ?", user).Find(&crosshairs)
	return crosshairs, tx.Error
}

func (p *psql) DeleteAllCrosshairsFromUser(user uuid.UUID) error {
	tx := p.db.Table(tableCrosshairs).Where("registrant_id = ?", user).Delete(&database.Crosshair{})
	return tx.Error
}

func (p *psql) DeleteCrosshairFromUserByCode(user uuid.UUID, crosshairCode string) error {
	tx := p.db.Table(tableCrosshairs).Where("registrant_id = ?", user).Where("code = ?", crosshairCode).Delete(&database.Crosshair{})
	return tx.Error
}

func (p *psql) EditCrosshairNote(ch *database.Crosshair) (*database.Crosshair, error) {
	tx := p.db.Table(tableCrosshairs).Where("registrant_id = ?", ch.RegistrantID).Where("code = ?", ch.Code).Update("note", ch.Note)
	return ch, tx.Error
}

func (p *psql) GetAllCrosshairsFromUserSortByDate(user uuid.UUID) ([]*database.Crosshair, error) {
	var crosshairs []*database.Crosshair
	tx := p.db.Table(tableCrosshairs).Order("created_at desc").Where("id = ?", user).Find(&crosshairs)
	return crosshairs, tx.Error
}
