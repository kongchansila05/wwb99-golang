package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"wwb99/config"
	"wwb99/models"

	"gorm.io/gorm"
)

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// GetPermissions handles GET /api/permissions or /api/roles/get?id=...
func GetPermissionRoles(w http.ResponseWriter, r *http.Request) {
	db := config.DB // Assuming config.DB is your *gorm.DB

	w.Header().Set("Content-Type", "application/json")

	// Check if a role ID is provided
	roleIDStr := r.URL.Query().Get("id")
	if roleIDStr != "" {
		// Fetch role with permissions
		roleID, err := strconv.Atoi(roleIDStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Status: false, Message: "Invalid role ID"})
			return
		}

		var role models.Role
		if err := db.Preload("Permissions").First(&role, roleID).Error; err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Response{Status: false, Message: "Role not found"})
			return
		}

		json.NewEncoder(w).Encode(Response{
			Status:  true,
			Message: "Role fetched successfully",
			Data:    role,
		})
		return
	}

	// Fetch all permissions
	var permissions []models.Permission
	if err := db.Find(&permissions).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: false, Message: "Error fetching permissions"})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Status:  true,
		Message: "Permissions fetched successfully",
		Data:    permissions,
	})
}
func GetRoleByID(w http.ResponseWriter, r *http.Request) {
	// Get "id" from query parameter
	idStr := strings.TrimSpace(r.URL.Query().Get("id"))
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	// Convert id to integer
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	// Fetch role with associated permissions
	var role models.Role
	if err := config.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Role not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error fetching role", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  true,
		"message": "Success",
		"data":    role,
	})
}

// GetRoles returns roles with permissions (paginated, searchable, sortable)
func GetRoles(w http.ResponseWriter, r *http.Request) {
	var roles []models.Role

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")
	sortField := r.URL.Query().Get("sortBy")
	order := r.URL.Query().Get("order")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	db := config.DB.Model(&models.Role{}).Preload("Permissions")

	// Search
	if search != "" {
		like := "%" + search + "%"
		db = db.Where("name LIKE ?", like)
	}

	var total int64
	db.Count(&total)

	// Sorting
	validSortFields := map[string]bool{
		"id":         true,
		"name":       true,
		"created_at": true,
	}
	if !validSortFields[sortField] {
		sortField = "created_at"
	}
	if strings.ToLower(order) != "asc" {
		order = "desc"
	}

	result := db.Order(sortField + " " + order).
		Limit(limit).
		Offset(offset).
		Find(&roles)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":       roles,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int((total + int64(limit) - 1) / int64(limit)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateRole creates a new role with permissions
func CreateRole(w http.ResponseWriter, r *http.Request) {
	var role models.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := config.DB.Create(&role).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	config.DB.Preload("Permissions").First(&role, role.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Role created successfully",
		"data":    role,
	})
}

// UpdateRole updates role name and permissions
func UpdateRole(w http.ResponseWriter, r *http.Request) {
	var role models.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if role.ID == 0 {
		http.Error(w, "Missing ID in request", http.StatusBadRequest)
		return
	}

	var existing models.Role
	if err := config.DB.Preload("Permissions").First(&existing, role.ID).Error; err != nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	// Update role name
	existing.Name = role.Name

	// Replace permissions
	if len(role.Permissions) > 0 {
		if err := config.DB.Model(&existing).Association("Permissions").Replace(role.Permissions); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := config.DB.Save(&existing).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	config.DB.Preload("Permissions").First(&existing, existing.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Role updated successfully",
		"data":    existing,
	})
}

// DeleteRole deletes a role by ID
func DeleteRole(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing role ID", http.StatusBadRequest)
		return
	}

	result := config.DB.Delete(&models.Role{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Role not found or already deleted", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Role deleted successfully"})
}
