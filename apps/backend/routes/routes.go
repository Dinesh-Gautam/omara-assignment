package routes

import (
	"strategic-insight-analyst/backend/handlers"
	"strategic-insight-analyst/backend/internal/api"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/auth/signup", handlers.SignupHandler).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(api.AuthMiddleware)
	protected.HandleFunc("/protected", handlers.ProtectedHandler).Methods("GET")
	protected.HandleFunc("/documents/upload", handlers.UploadDocumentHandler).Methods("POST")
	protected.HandleFunc("/documents", handlers.GetDocumentsHandler).Methods("GET")
	protected.HandleFunc("/documents/download/{document_id}", handlers.DownloadDocumentHandler).Methods("GET")
	protected.HandleFunc("/documents/{document_id}/status", handlers.GetDocumentStatusHandler).Methods("GET")
	protected.HandleFunc("/documents/{document_id}", handlers.DeleteDocumentHandler).Methods("DELETE")
	protected.HandleFunc("/chat", handlers.ChatHandler).Methods("POST")
	protected.HandleFunc("/chat/{document_id}", handlers.GetChatHistoryHandler).Methods("GET")
}
