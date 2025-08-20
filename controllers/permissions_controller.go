package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"wwb99/config"
	"wwb99/models"
)

// Request body structure
type AssignPermissionsRequest struct {
	ID          uint   `json:"id"`          // role ID
	Permissions []uint `json:"permissions"` // e.g. [39, 40, 41]
}

func AssignPermissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := config.DB

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request body
	var req AssignPermissionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.ID == 0 {
		http.Error(w, "Missing role ID", http.StatusBadRequest)
		return
	}

	if len(req.Permissions) == 0 {
		http.Error(w, "No permissions provided", http.StatusBadRequest)
		return
	}

	// Check if role exists
	var role models.Role
	if err := db.First(&role, req.ID).Error; err != nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	// Start transaction
	tx := db.Begin()
	if tx.Error != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	// Delete old permissions for this role
	if err := tx.Where("role_id = ?", req.ID).Delete(&models.RolePermission{}).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to delete old permissions", http.StatusInternalServerError)
		return
	}
	assignedIDs := req.Permissions
	for _, pid := range assignedIDs {
		rp := models.RolePermission{
			RoleID:       req.ID,
			PermissionID: pid,
		}
		if err := tx.Create(&rp).Error; err != nil {
			tx.Rollback()
			http.Error(w, fmt.Sprintf("Failed to assign permission ID %d: %v", pid, err), http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()

	// Return success response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      true,
		"message":     "Permissions assigned successfully",
		"roleId":      role.ID,
		"assignedIDs": req.Permissions,
	})
}

// GetPermissions returns paginated permissions with search/sorting
func GetPermissions(w http.ResponseWriter, r *http.Request) {
	var permissions []models.Permission

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")
	sortField := r.URL.Query().Get("sortBy")
	order := r.URL.Query().Get("order")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		http.Error(w, "'page' must be a positive integer", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		http.Error(w, "'limit' must be a positive integer", http.StatusBadRequest)
		return
	}

	offset := (page - 1) * limit
	db := config.DB.Model(&models.Permission{})

	// Search by name
	if search != "" {
		like := "%" + search + "%"
		db = db.Where("name LIKE ?", like)
	}

	var total int64
	db.Count(&total)

	// Valid sort fields
	validSortFields := map[string]bool{
		"id":         true,
		"name":       true,
		"created_at": true,
	}

	if !validSortFields[sortField] {
		sortField = "created_at"
	}

	order = strings.ToLower(order)
	if order != "asc" {
		order = "desc"
	}

	result := db.Order(sortField + " " + order).
		Limit(limit).
		Offset(offset).
		Find(&permissions)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":       permissions,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int((total + int64(limit) - 1) / int64(limit)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreatePermission creates a new permission
func CreatePermission(w http.ResponseWriter, r *http.Request) {
	var permission models.Permission
	if err := json.NewDecoder(r.Body).Decode(&permission); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := config.DB.Create(&permission).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Permission created successfully",
		"data":    permission,
	})
}

// UpdatePermission updates an existing permission
func UpdatePermission(w http.ResponseWriter, r *http.Request) {
	var permission models.Permission

	if err := json.NewDecoder(r.Body).Decode(&permission); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if permission.ID == 0 {
		http.Error(w, "Missing ID in request", http.StatusBadRequest)
		return
	}

	var existing models.Permission
	if err := config.DB.First(&existing, permission.ID).Error; err != nil {
		http.Error(w, "Permission not found", http.StatusNotFound)
		return
	}

	err := config.DB.Model(&existing).Updates(models.Permission{
		Name: permission.Name,
	}).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Permission updated successfully"})
}

// DeletePermission deletes a permission by ID
func DeletePermission(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing permission ID", http.StatusBadRequest)
		return
	}

	result := config.DB.Delete(&models.Permission{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Permission not found or already deleted", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Permission deleted successfully"})
}
