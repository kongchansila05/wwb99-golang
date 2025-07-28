package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
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
		"data": news,
	})
}
func GetNews(w http.ResponseWriter, r *http.Request) {
	var newsList []models.News

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
	config.DB.Model(&models.News{}).Count(&total)

	// Fetch paginated results
	result := config.DB.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&newsList)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Response JSON structure
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

	// Set created_at (if not using GORM's auto time tracking)
	// news.CreatedAt = time.Now() // optional if your model uses `gorm:"autoCreateTime"`

	// Save using GORM
	if err := config.DB.Create(&news).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the created object as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}
func UpdateNews(w http.ResponseWriter, r *http.Request) {
	var news models.News

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure ID is provided
	if news.ID == 0 {
		http.Error(w, "Missing ID in request", http.StatusBadRequest)
		return
	}

	// Check if the news record exists
	var existing models.News
	if err := config.DB.First(&existing, news.ID).Error; err != nil {
		http.Error(w, "News not found", http.StatusNotFound)
		return
	}

	// Perform the update
	err := config.DB.Model(&existing).Updates(models.News{
		Title:     news.Title,
		Image:     news.Image,
		Detail:    news.Detail,
		CreatedBy: news.CreatedBy,
	}).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "News updated successfully"})
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
