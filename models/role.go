package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string       `gorm:"unique"`
	Permissions []Permission `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
