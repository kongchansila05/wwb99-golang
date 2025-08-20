package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"wwb99/config"
	"wwb99/models"
)

func GetNewsHome(w http.ResponseWriter, r *http.Request) {
	var newsList []models.News
	result := config.DB.Order("created_at DESC").Limit(4).Find(&newsList)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newsList)
}

func GetNewsByID(w http.ResponseWriter, r *http.Request) {
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

	var news models.News
	result := config.DB.First(&news, id)
	if result.Error != nil {
		http.Error(w, "News not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Success",
		"data":    news,
	})
}

func GetNews(w http.ResponseWriter, r *http.Request) {
	var newsList []models.News

	// Parse query params
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")
	sortField := r.URL.Query().Get("sortBy")
	order := r.URL.Query().Get("order")

	// Parse page and limit, error if invalid or missing
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

	db := config.DB.Model(&models.News{})

	// Apply search filter if present
	if search != "" {
		likeQuery := "%" + search + "%"
		db = db.Where("title LIKE ? OR detail LIKE ?", likeQuery, likeQuery)
	}

	// Count total records after search filter
	var total int64
	db.Count(&total)

	// Validate sortField, allow new fields: id, title, created_at, created_by
	validSortFields := map[string]bool{
		"id":         true,
		"title":      true,
		"created_at": true,
		"created_by": true,
	}

	if !validSortFields[sortField] {
		sortField = "created_at" // default sort field
	}

	// Validate order param
	order = strings.ToLower(order)
	if order != "asc" {
		order = "desc" // default order
	}

	// Apply sorting, limit and offset
	result := db.Order(sortField + " " + order).
		Limit(limit).
		Offset(offset).
		Find(&newsList)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Build and send response
	response := map[string]interface{}{
		"data":       newsList,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": int((total + int64(limit) - 1) / int64(limit)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateNews(w http.ResponseWriter, r *http.Request) {
	var news models.News
	// Decode JSON body into news struct
	if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Save using GORM
	if err := config.DB.Create(&news).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Prepare response
	response := struct {
		Message string      `json:"message"`
		Data    models.News `json:"data"`
	}{
		Message: "News created successfully",
		Data:    news,
	}
	// Return the created object as JSON with message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// ✅ Get ID from query parameter
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"message":"Missing id parameter"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, `{"message":"Invalid id parameter"}`, http.StatusBadRequest)
		return
	}

	// ✅ Decode JSON body into struct
	var updatedData models.News
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, `{"message":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// ✅ Check if the record exists
	var existing models.News
	if err := config.DB.First(&existing, id).Error; err != nil {
		http.Error(w, `{"message":"News not found"}`, http.StatusNotFound)
		return
	}

	// ✅ Update fields
	existing.Title = updatedData.Title
	existing.Image = updatedData.Image
	existing.Detail = updatedData.Detail
	existing.Content = updatedData.Content
	existing.CreatedBy = updatedData.CreatedBy

	// ✅ Save to DB
	if err := config.DB.Save(&existing).Error; err != nil {
		http.Error(w, `{"message":"Failed to update news"}`, http.StatusInternalServerError)
		return
	}

	// ✅ JSON Response
	response := struct {
		Message string      `json:"message"`
		Data    models.News `json:"data"`
	}{
		Message: "News updated successfully",
		Data:    existing,
	}

	json.NewEncoder(w).Encode(response)
}

func DeleteNews(w http.ResponseWriter, r *http.Request) {
	// Parse the news ID from the request query parameters
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing news ID", http.StatusBadRequest)
		return
	}

	// Attempt to delete the news with the given ID
	result := config.DB.Delete(&models.News{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "News not found or already deleted", http.StatusNotFound)
		return
	}

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "News deleted successfully"})
}
