package main

import (
	"log"
	"net/http"

	"strategic-insight-analyst/backend/config"
	"strategic-insight-analyst/backend/database"
	"strategic-insight-analyst/backend/firebase"
	"strategic-insight-analyst/backend/internal/storage"
	"strategic-insight-analyst/backend/routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	database.Connect()
	database.Migrate()
	firebase.Initialize()
	if err := storage.InitializeGCS(); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	// CORS configuration
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{config.AppConfig.FrontendURL}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
	)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", corsHandler(r)); err != nil {
		log.Fatal(err)
	}
}
