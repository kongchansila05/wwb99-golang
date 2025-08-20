package routes

import (
	"wwb99/controllers"
	"wwb99/middleware"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	// start admin
	r.HandleFunc("/api/register", controllers.Register).Methods("POST")
	r.HandleFunc("/api/login", controllers.Login).Methods("POST")
	r.HandleFunc("/api/refresh", controllers.RefreshToken).Methods("POST")

	r.HandleFunc("/api/news", controllers.GetNews).Methods("GET")
	r.HandleFunc("/api/news/create", controllers.CreateNews).Methods("POST")
	r.HandleFunc("/api/news/update/{id}", controllers.UpdateNews).Methods("PUT")
	r.HandleFunc("/api/news/delete", controllers.DeleteNews)
	r.HandleFunc("/api/news/getbyid", controllers.GetNewsByID)

	r.HandleFunc("/api/highlights", controllers.GetHighlights).Methods("GET")
	r.HandleFunc("/api/highlights/create", controllers.CreateHighlights).Methods("POST")
	r.HandleFunc("/api/highlights/update", controllers.UpdateHighlights).Methods("PUT")
	r.HandleFunc("/api/highlights/delete", controllers.DeleteHighlights)
	r.HandleFunc("/api/highlights/getbyid", controllers.GetHighlightsByID)

	r.HandleFunc("/api/footers", controllers.GetFooters).Methods("GET")
	r.HandleFunc("/api/footers/create", controllers.CreateFooter).Methods("POST")
	r.HandleFunc("/api/footers/update", controllers.UpdateFooter).Methods("PUT")
	r.HandleFunc("/api/footers/delete", controllers.DeleteFooter)
	r.HandleFunc("/api/footers/getbyid", controllers.GetFooterByID)

	r.HandleFunc("/api/sponsors", controllers.GetSponsors).Methods("GET")
	r.HandleFunc("/api/sponsors/create", controllers.CreateSponsor).Methods("POST")
	r.HandleFunc("/api/sponsors/update", controllers.UpdateSponsor).Methods("PUT")
	r.HandleFunc("/api/sponsors/delete", controllers.DeleteSponsor)
	r.HandleFunc("/api/sponsors/getbyid", controllers.GetSponsorByID)

	r.HandleFunc("/api/permissions", controllers.GetPermissions).Methods("GET")
	r.HandleFunc("/api/permissions/create", controllers.CreatePermission).Methods("POST")
	r.HandleFunc("/api/permissions/update", controllers.UpdatePermission).Methods("PUT")

	r.HandleFunc("/api/roles", controllers.GetRoles).Methods("GET")
	r.HandleFunc("/api/roles", controllers.CreateRole).Methods("POST")
	r.HandleFunc("/api/roles", controllers.UpdateRole).Methods("PUT")
	r.HandleFunc("/api/roles", controllers.DeleteRole).Methods("DELETE")
	r.HandleFunc("/api/roles/getbyid", controllers.GetRoleByID).Methods("GET")
	r.HandleFunc("/api/roles/permissions", controllers.GetPermissionRoles).Methods("GET")
	r.HandleFunc("/api/roles/assign", controllers.AssignPermissions).Methods("PUT")

	// end admin

	// start client
	r.HandleFunc("/api/news_home", controllers.GetNewsHome).Methods("GET")
	r.HandleFunc("/api/highlights_home", controllers.GetHighlightsHome).Methods("GET")
	r.HandleFunc("/api/footers_home", controllers.GetFootersHome).Methods("GET")
	r.HandleFunc("/api/sponsors_home", controllers.GetSponsorsHome).Methods("GET")

	// end client

	secured := r.PathPrefix("/api").Subrouter()
	secured.Use(middleware.AuthMiddleware)
	secured.HandleFunc("/profile", controllers.Profile).Methods("GET")

	return r
}
