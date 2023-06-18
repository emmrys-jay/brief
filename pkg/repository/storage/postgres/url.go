package postgres

import (
	"brief/internal/model"
	"context"
)

// CreateURL stores 'url' in the database
func (p *Postgres) CreateURL(ctx context.Context, url *model.URL) error {
	db, cancel := p.DBWithTimeout(ctx)
	defer cancel()
	return db.Create(url).Error
}

// GetURL fetches a url entry from the database using its 'hash'
func (p *Postgres) GetURL(ctx context.Context, hash string) (*model.URL, error) {
	db, cancel := p.DBWithTimeout(ctx)
	defer cancel()

	var url model.URL
	err := db.First(&url, "hash = ?", hash).Error
	return &url, err
}
