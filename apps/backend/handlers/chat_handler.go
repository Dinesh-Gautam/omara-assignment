package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"strategic-insight-analyst/backend/services"
	"strategic-insight-analyst/backend/utils"

	"firebase.google.com/go/auth"
	"github.com/gorilla/mux"
)

type chatRequest struct {
	DocumentID        string   `json:"document_id"`
	UserMessage       string   `json:"message"`
	AttachedDocuments []string `json:"attached_documents"`
}

type streamData struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*auth.Token)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "User not found in context")
		return
	}
	userID := user.UID

	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get attached document details
	attachedDocs, err := services.GetDocumentsByIDs(req.AttachedDocuments, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve attached documents: "+err.Error())
		return
	}

	if _, err := services.SaveUserMessage(req.DocumentID, userID, req.UserMessage, attachedDocs); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save user message: "+err.Error())
		return
	}

	contextText, err := services.GetRelevantContext(req.DocumentID, req.UserMessage, req.AttachedDocuments, userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get relevant context: "+err.Error())
		return
	}

	history, err := services.GetChatHistoryForLLM(req.DocumentID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get chat history for LLM: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "Streaming unsupported!")
		return
	}

	streamChan := make(chan string)
	var fullResponse strings.Builder
	var streamErr error
	done := make(chan bool)

	go func() {
		defer close(done)
		for chunk := range streamChan {
			fullResponse.WriteString(chunk)
			jsonData, err := json.Marshal(streamData{Token: chunk})
			if err != nil {
				streamErr = fmt.Errorf("error marshalling stream data: %w", err)
				return
			}
			_, err = fmt.Fprintf(w, "data: %s\n\n", jsonData)
			if err != nil {
				streamErr = err
				return
			}
			flusher.Flush()
		}
	}()

	// The service function now returns the full response and a potential error
	aiResponse, serviceErr := services.StreamChatResponse(req.UserMessage, contextText, history, streamChan)
	<-done

	if serviceErr != nil {
		log.Printf("Error from streaming service: %v", serviceErr)
		// Send an error message to the client through the stream
		errorData, _ := json.Marshal(streamData{Error: "Failed to get response from AI. Please try again."})
		fmt.Fprintf(w, "data: %s\n\n", errorData)
		flusher.Flush()
		return
	}

	if streamErr != nil {
		log.Printf("Error writing to stream: %v", streamErr)
		// Don't try to write to the stream if it's already broken
		return
	}

	// Save the successful response
	if _, err := services.SaveAIMessage(req.DocumentID, userID, aiResponse); err != nil {
		log.Printf("Failed to save AI response: %v", err)
	}
}

func GetChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	documentID := vars["document_id"]

	history, err := services.GetChatHistory(documentID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve chat history: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, history)
}
