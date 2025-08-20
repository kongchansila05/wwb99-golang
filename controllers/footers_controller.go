package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"wwb99/config"
	"wwb99/models"
)

func GetFootersHome(w http.ResponseWriter, r *http.Request) {
	var footers []models.Footers
	result := config.DB.Order("created_at DESC").Find(&footers)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Success",
		"data":    footers,
	})
}

// Get all footers with pagination, search, sorting
func GetFooters(w http.ResponseWriter, r *http.Request) {
	var footers []models.Footers

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")
	sortField := r.URL.Query().Get("sortBy")
	order := r.URL.Query().Get("order")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		http.Error(w, "'page' query parameter is required and must be a positive integer", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		http.Error(w, "'limit' query parameter is required and must be a positive integer", http.StatusBadRequest)
		return
	}

	offset := (page - 1) * limit
	db := config.DB.Model(&models.Footers{})

	// Search by name or redirect
	if search != "" {
		likeQuery := "%" + search + "%"
		db = db.Where("name LIKE ? OR redirect LIKE ?", likeQuery, likeQuery)
	}

	var total int64
	db.Count(&total)

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
		Find(&footers)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":       footers,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int((total + int64(limit) - 1) / int64(limit)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Get footer by ID
func GetFooterByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	var footer models.Footers
	result := config.DB.First(&footer, id)
	if result.Error != nil {
		http.Error(w, "Footer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Success",
		"data":    footer,
	})
}

// Create new footer
func CreateFooter(w http.ResponseWriter, r *http.Request) {
	var footer models.Footers
	if err := json.NewDecoder(r.Body).Decode(&footer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := config.DB.Create(&footer).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string         `json:"message"`
		Data    models.Footers `json:"data"`
	}{
		Message: "Footer created successfully",
		Data:    footer,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func UpdateFooter(w http.ResponseWriter, r *http.Request) {
	var Footers models.Footers

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&Footers); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure ID is provided
	if Footers.ID == 0 {
		http.Error(w, "Missing ID in request", http.StatusBadRequest)
		return
	}

	// Check if the Footers record exists
	var existing models.Footers
	if err := config.DB.First(&existing, Footers.ID).Error; err != nil {
		http.Error(w, "Footers not found", http.StatusNotFound)
		return
	}

	// Perform the update
	err := config.DB.Model(&existing).Updates(models.Footers{
		Name:     Footers.Name,
		ImageURL: Footers.ImageURL,
		Redirect: Footers.Redirect,
	}).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Footers updated successfully"})
}

// Delete footer
func DeleteFooter(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing footer ID", http.StatusBadRequest)
		return
	}

	result := config.DB.Delete(&models.Footers{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Footer not found or already deleted", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Footer deleted successfully"})
}
