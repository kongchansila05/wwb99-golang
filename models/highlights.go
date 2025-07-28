package models

import "time"

type Highlights struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Image     string    `json:"image"`
	Detail    string    `json:"detail"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
