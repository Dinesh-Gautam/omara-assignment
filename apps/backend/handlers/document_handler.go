package handlers

import (
	"fmt"
	"net/http"

	"strategic-insight-analyst/backend/services"
	"strategic-insight-analyst/backend/utils"

	"firebase.google.com/go/auth"
	"github.com/gorilla/mux"
)

func UploadDocumentHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*auth.Token)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "User not found in context")
		return
	}
	userID := user.UID

	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		utils.RespondWithError(w, http.StatusBadRequest, "File size exceeds 10MB limit")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Could not retrieve file from form")
		return
	}
	defer file.Close()

	// Validate file type
	contentType := handler.Header.Get("Content-Type")
	if contentType != "application/pdf" && contentType != "text/plain" {
		utils.RespondWithError(w, http.StatusUnsupportedMediaType, "Invalid file type: only application/pdf and text/plain are allowed")
		return
	}

	doc, err := services.ProcessAndSaveDocument(file, handler, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process and save document: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, doc)
}

func GetDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*auth.Token)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	userID := user.UID

	documents, err := services.GetUserDocuments(userID)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve documents: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, documents)
}

func GetDocumentStatusHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*auth.Token)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "User not found in context")
		return
	}
	userID := user.UID

	vars := mux.Vars(r)
	documentID := vars["document_id"]

	doc, err := services.GetDocumentStatus(documentID, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get document status: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, doc)
}

func DownloadDocumentHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*auth.Token)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "User not found in context")
		return
	}
	userID := user.UID

	vars := mux.Vars(r)
	documentID := vars["document_id"]

	fileContent, fileName, err := services.DownloadDocument(documentID, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to download document: "+err.Error())
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileContent)
}

func DeleteDocumentHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*auth.Token)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "User not found in context")
		return
	}
	userID := user.UID

	vars := mux.Vars(r)
	documentID := vars["document_id"]

	err := services.DeleteDocument(documentID, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete document: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Document deleted successfully"})
}
