package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strategic-insight-analyst/backend/database"
	"strategic-insight-analyst/backend/firebase"
	"strategic-insight-analyst/backend/models"
	"strategic-insight-analyst/backend/utils"
	"time"
)

type tokenRequest struct {
	Token string `json:"token"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var req tokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, err := firebase.AuthClient.VerifyIDToken(context.Background(), req.Token)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid Firebase token: "+err.Error())
		return
	}

	user, err := database.FindUserByID(token.UID)
	if err == sql.ErrNoRows {
		firebaseUser, err := firebase.AuthClient.GetUser(context.Background(), token.UID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get user from Firebase: "+err.Error())
			return
		}

		log.Println(firebaseUser)

		authMethod := "google"
		if firebaseUser.Email == "" {
			authMethod = "guest"
		}

		newUser := &models.User{
			ID:         token.UID,
			Email:      firebaseUser.Email,
			AuthMethod: authMethod,
			CreatedAt:  time.Now(),
		}

		if authMethod == "guest" {
			newUser.Email = "guest_" + token.UID + "@example.com"
		}

		if err := database.CreateUser(newUser); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user in database: "+err.Error())
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, newUser)
		return
	} else if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to find user: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, user)
}
