package handlers

import (
	"net/http"
	"strategic-insight-analyst/backend/utils"
)

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "This is a protected route"})
}
