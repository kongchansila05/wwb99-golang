package main

import (
	"log"
	"net/http"
	"os"

	"wwb99/config"
	"wwb99/middleware"
	"wwb99/models"
	"wwb99/routes"
)

func main() {
	// Connect and migrate DB
	config.Connect()
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
	)
	if err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}

	// 🔁 Load your app router
	router := routes.RegisterRoutes()

	// ✅ Wrap with prerender first, then CORS (order matters)
	withPrerender := middleware.PrerenderMiddleware(router)
	withCORS := middleware.CORSMiddleware(withPrerender)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Server running at http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, withCORS); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
