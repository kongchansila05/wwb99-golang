package seeder

import (
	"log"
	"wwb99/config"
	"wwb99/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedRolesAndPermissions() {
	db := config.DB

	// 1. Create or get permissions
	permNames := []string{"view_users", "edit_users", "delete_users", "view_roles", "edit_roles", "delete_roles", "view_permissions", "edit_permissions", "delete_permissions"}
	var permissions []models.Permission

	for _, name := range permNames {
		var p models.Permission
		db.FirstOrCreate(&p, models.Permission{Name: name})
		permissions = append(permissions, p)
	}

	// 2. Create or get roles
	var adminRole models.Role
	db.FirstOrCreate(&adminRole, models.Role{Name: "admin"})

	var userRole models.Role
	db.FirstOrCreate(&userRole, models.Role{Name: "user"})

	// 3. Attach permissions to adminRole
	err := db.Model(&adminRole).Association("Permissions").Replace(permissions)
	if err != nil {
		log.Printf("Error attaching permissions: %v", err)
	}

	log.Println("✅ Seeded roles and permissions.")
}
func SeedOwnerUser() {
	db := config.DB

	// 1. Create role "owner" if not exists
	var ownerRole models.Role
	db.FirstOrCreate(&ownerRole, models.Role{Name: "owner"})

	// 2. Create user with that role
	password, _ := bcrypt.GenerateFromPassword([]byte("owner123"), bcrypt.DefaultCost)

	var user models.User
	db.FirstOrCreate(&user, models.User{
		Username: "owner",
		Password: string(password),
		RoleID:   ownerRole.ID,
	})

	log.Println("✅ Seeded owner user with role 'owner'")
}
