package models

import "time"

type Footers struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	ImageURL  string    `json:"image_url" gorm:"type:varchar(512)"`
	Redirect  string    `json:"redirect" gorm:"type:varchar(512)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
