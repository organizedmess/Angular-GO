package models

import "time"

type URL struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OriginalURL string    `gorm:"type:text;not null" json:"original_url"`
	ShortCode   string    `gorm:"type:varchar(10);uniqueIndex;not null" json:"short_code"`
	ClickCount  int       `gorm:"default:0" json:"click_count"`
	CreatedAt   time.Time `json:"created_at"`
}
