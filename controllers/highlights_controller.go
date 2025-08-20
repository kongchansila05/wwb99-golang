package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wwb99/config"
	"wwb99/models"
)

func GetHighlightsHome(w http.ResponseWriter, r *http.Request) {
	var highlightsList []models.Highlights
	result := config.DB.Order("created_at DESC").Limit(4).Find(&highlightsList)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(highlightsList)
}
func GetHighlightsByID(w http.ResponseWriter, r *http.Request) {
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

	var highlights models.Highlights
	result := config.DB.First(&highlights, id)
	if result.Error != nil {
		http.Error(w, "Highlights not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Success",
		"data":    highlights,
	})
}
func GetHighlights(w http.ResponseWriter, r *http.Request) {
	var highlightsList []models.Highlights

	// Parse pagination query params
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	offset := (page - 1) * limit

	// Get total record count
	var total int64
	config.DB.Model(&models.Highlights{}).Count(&total)

	// Fetch paginated results
	result := config.DB.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&highlightsList)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Response JSON structure
	response := map[string]interface{}{
		"data":       highlightsList,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int((total + int64(limit) - 1) / int64(limit)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func CreateHighlights(w http.ResponseWriter, r *http.Request) {
	var Highlights models.Highlights

	// Decode JSON body into Highlights struct
	if err := json.NewDecoder(r.Body).Decode(&Highlights); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set created_at (if not using GORM's auto time tracking)
	// Highlights.CreatedAt = time.Now() // optional if your model uses `gorm:"autoCreateTime"`

	// Save using GORM
	if err := config.DB.Create(&Highlights).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the created object as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Highlights)
}
func UpdateHighlights(w http.ResponseWriter, r *http.Request) {
	var Highlights models.Highlights

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&Highlights); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure ID is provided
	if Highlights.ID == 0 {
		http.Error(w, "Missing ID in request", http.StatusBadRequest)
		return
	}

	// Check if the Highlights record exists
	var existing models.Highlights
	if err := config.DB.First(&existing, Highlights.ID).Error; err != nil {
		http.Error(w, "Highlights not found", http.StatusNotFound)
		return
	}

	// Perform the update
	err := config.DB.Model(&existing).Updates(models.Highlights{
		Title:     Highlights.Title,
		Image:     Highlights.Image,
		Content:   Highlights.Content,
		CreatedBy: Highlights.CreatedBy,
	}).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Highlights updated successfully"})
}
func DeleteHighlights(w http.ResponseWriter, r *http.Request) {
	// Parse the Highlights ID from the request query parameters
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing Highlights ID", http.StatusBadRequest)
		return
	}

	// Attempt to delete the Highlights with the given ID
	result := config.DB.Delete(&models.Highlights{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Highlights not found or already deleted", http.StatusNotFound)
		return
	}

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Highlights deleted successfully"})
}
