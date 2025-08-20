package models

type RolePermission struct {
	RoleID       uint `gorm:"column:role_id"`
	PermissionID uint `gorm:"column:permission_id"`
}

// Optional: explicitly set table name
func (RolePermission) TableName() string {
	return "role_permissions"
}
