package controllers

import (
	"encoding/json"
	"net/http"
	"wwb99/config"
	"wwb99/models"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uint)

	var user models.User
	config.DB.Preload("Role.Permissions").First(&user, userID)

	json.NewEncoder(w).Encode(user)
}
