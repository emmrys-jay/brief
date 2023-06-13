package model

import "time"

type URL struct {
	ID        string    `json:"id,omitempty" gorm:"column:id;index;unique;not null;type:varchar(50)"`
	LongURL   string    `json:"long_url,omitempty" gorm:"column:long_url;not null" validate:"required"`
	ShortURL  string    `json:"short_url,omitempty" gorm:"column:short_url;not null;unique;index"`
	UserID    string    `json:"user_id,omitempty" gorm:"column:user_id;type:varchar(50)"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;index"`
}
