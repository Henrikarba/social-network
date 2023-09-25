package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"social-network/pkg/models"

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
}
