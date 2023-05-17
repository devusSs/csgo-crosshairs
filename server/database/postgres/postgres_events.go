package postgres

import "github.com/devusSs/crosshairs/database"

func (p *psql) AddEvent(event *database.Event) (*database.Event, error) {
	tx := p.db.Table(tableEvents).Create(&event)
	return event, tx.Error
}

func (p *psql) GetEvents() ([]*database.Event, error) {
	var events []*database.Event
	tx := p.db.Table(tableEvents).Find(&events)
	return events, tx.Error
}

func (p *psql) GetEventsByType(eventType string) ([]*database.Event, error) {
	var events []*database.Event
	tx := p.db.Table(tableEvents).Where("type = ?", eventType).Find(&events)
	return events, tx.Error
}
