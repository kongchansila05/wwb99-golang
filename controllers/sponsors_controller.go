package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"wwb99/config"
	"wwb99/models"
)

func GetSponsorsHome(w http.ResponseWriter, r *http.Request) {
	var sponsors []models.Sponsors
	result := config.DB.Order("created_at DESC").Find(&sponsors)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Success",
		"data":    sponsors,
	})
}

// Get all sponsors with pagination, search, sorting
func GetSponsors(w http.ResponseWriter, r *http.Request) {
	var sponsors []models.Sponsors

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
	db := config.DB.Model(&models.Sponsors{})

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
		Find(&sponsors)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":       sponsors,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int((total + int64(limit) - 1) / int64(limit)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Get sponsor by ID
func GetSponsorByID(w http.ResponseWriter, r *http.Request) {
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

	var sponsor models.Sponsors
	result := config.DB.First(&sponsor, id)
	if result.Error != nil {
		http.Error(w, "Sponsor not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Success",
		"data":    sponsor,
	})
}

// Create new sponsor
func CreateSponsor(w http.ResponseWriter, r *http.Request) {
	var sponsor models.Sponsors
	if err := json.NewDecoder(r.Body).Decode(&sponsor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := config.DB.Create(&sponsor).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string          `json:"message"`
		Data    models.Sponsors `json:"data"`
	}{
		Message: "Sponsor created successfully",
		Data:    sponsor,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateSponsor(w http.ResponseWriter, r *http.Request) {
	var sponsor models.Sponsors

	if err := json.NewDecoder(r.Body).Decode(&sponsor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if sponsor.ID == 0 {
		http.Error(w, "Missing ID in request", http.StatusBadRequest)
		return
	}

	var existing models.Sponsors
	if err := config.DB.First(&existing, sponsor.ID).Error; err != nil {
		http.Error(w, "Sponsor not found", http.StatusNotFound)
		return
	}

	err := config.DB.Model(&existing).Updates(models.Sponsors{
		Name:     sponsor.Name,
		ImageURL: sponsor.ImageURL,
		Redirect: sponsor.Redirect,
	}).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Sponsor updated successfully"})
}

// Delete sponsor
func DeleteSponsor(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing sponsor ID", http.StatusBadRequest)
		return
	}

	result := config.DB.Delete(&models.Sponsors{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Sponsor not found or already deleted", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Sponsor deleted successfully"})
}
