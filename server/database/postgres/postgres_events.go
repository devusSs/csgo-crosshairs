package postgres

import "github.com/devusSs/crosshairs/database"

func (p *psql) AddEvent(event *database.Event) (*database.Event, error) {
	tx := p.db.Table(tableEvents).Create(&event)
	return event, tx.Error
}

func (p *psql) GetEvents() ([]*database.Event, error) {
	var events []*database.Event
	tx := p.db.Table(tableEvents).Order("created_at desc").Find(&events)
	return events, tx.Error
}

func (p *psql) GetEventsByType(eventType string) ([]*database.Event, error) {
	var events []*database.Event
	tx := p.db.Table(tableEvents).Order("created_at desc").Where("type = ?", eventType).Find(&events)
	return events, tx.Error
}

func (p *psql) GetEventsWithLimit(limit int) ([]*database.Event, error) {
	var events []*database.Event
	tx := p.db.Table(tableEvents).Order("created_at desc").Limit(limit).Find(&events)
	return events, tx.Error
}

func (p *psql) GetEventsByTypeWithLimit(eventType string, limit int) ([]*database.Event, error) {
	var events []*database.Event
	tx := p.db.Table(tableEvents).Order("created_at desc").Limit(limit).Where("type = ?", eventType).Find(&events)
	return events, tx.Error
}
