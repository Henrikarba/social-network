package server

import (
	"encoding/json"
	"net/http"

	"social-network/pkg/models"
	"social-network/pkg/utils"
)

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var login models.UserRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&login); err != nil {
		http.Error(w, "Unable to decode JSON data", http.StatusBadRequest)
		return
	}
	valid, id, _ := models.ValidateLogin(s.db.DB, login.Email, login.Password)
	if !valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tkn, uuid, _ := utils.GenerateJWT(id)

	newSession := &models.Session{
		UserID: id,
		UUID:   uuid,
	}

	err := models.NewSession(s.db.DB, *newSession)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    tkn,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	response := struct {
		Message string `json:"message"`
	}{
		Message: "Login successful",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	// Set Content-Type header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
