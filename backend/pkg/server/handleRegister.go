package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"social-network/pkg/models"
	"social-network/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var register models.UserRequest
	var MimeType string
	var hashPassword string

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&register); err != nil {
		http.Error(w, "Unable to decode JSON data", http.StatusBadRequest)
		return
	}
	if register.FirstName == "" {
		http.Error(w, "First Name required", http.StatusBadRequest)
		return
	}

	if register.LastName == "" {
		http.Error(w, "Last Name required", http.StatusBadRequest)
		return
	}

	if register.Email == "" {
		http.Error(w, "Email required", http.StatusBadRequest)
		return
	} else {
		isValid := isValidEmail(register.Email)
		if !isValid {
			http.Error(w, "Email not valid", http.StatusBadRequest)
			return
		}
	}

	if register.Password == "" {
		http.Error(w, "Password required", http.StatusBadRequest)
		return
	}

	if register.DateOfBirth == "" {
		http.Error(w, "Date Of Birth required", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err == nil {
		MimeType = header.Header.Get("Content-Type")
		register.ImageData, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read file data", http.StatusInternalServerError)
			return
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	hashPassword = string(hashedPassword)

	err = models.RegisterUser(s.db.DB, register.Email, hashPassword, register.FirstName, register.LastName, register.DateOfBirth, register.Nickname, register.AboutMe, MimeType, register.ImageData)
	if err != nil {
		fmt.Println(err)
	}

	// auto login after register

	valid, id, _ := models.ValidateLogin(s.db.DB, register.Email, register.Password)
	if !valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tkn, uuid, _ := utils.GenerateJWT(id)

	newSession := &models.Session{
		UserID: id,
		UUID:   uuid,
	}

	err = models.NewSession(s.db.DB, *newSession)
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

func isValidEmail(email string) bool {
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailPattern)

	return re.MatchString(email)
}
