package controllers

import (
	"encoding/json"
	"net/http"
	"wwb99/config"
	"wwb99/models"
	"wwb99/utils"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashed)

	// Assign default role (ID = 1)
	user.RoleID = 1

	config.DB.Create(&user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input models.User
	json.NewDecoder(r.Body).Decode(&input)

	var user models.User
	config.DB.Preload("Role.Permissions").Where("username = ?", input.Username).First(&user)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, _ := utils.GenerateAccessToken(user.ID)
	refreshToken, _ := utils.GenerateRefreshToken(user.ID)

	response := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role.Name,
			"permissions": func() []string {
				perms := []string{}
				for _, p := range user.Role.Permissions {
					perms = append(perms, p.Name)
				}
				return perms
			}(),
		},
	}

	json.NewEncoder(w).Encode(response)
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var data struct {
		RefreshToken string `json:"refresh_token"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil || data.RefreshToken == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	claims, err := utils.ValidateRefreshToken(data.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	userID := uint(claims["user_id"].(float64))
	newAccessToken, _ := utils.GenerateAccessToken(userID)

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": newAccessToken,
	})
}
