package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"social-network/pkg/models"
	"social-network/pkg/utils"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}
	var postData models.UserRequest
	postData.Email = r.FormValue("email")
	postData.Nickname = r.FormValue("nickname")
	postData.FirstName = r.FormValue("first_name")
	postData.LastName = r.FormValue("last_name")
	postData.Password = r.FormValue("password")
	postData.DateOfBirth = r.FormValue("date_of_birth")
	file, header, err := r.FormFile("image")
	if err == nil {
		postData.MimeType = header.Header.Get("Content-Type")
		postData.ImageData, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read file data", http.StatusInternalServerError)
			return
		}
	}
	if postData.FirstName == "" {
		http.Error(w, "First Name required", http.StatusBadRequest)
		return
	}

	if postData.LastName == "" {
		http.Error(w, "Last Name required", http.StatusBadRequest)
		return
	}

	if postData.Email == "" {
		http.Error(w, "Email required", http.StatusBadRequest)
		return
	} else {
		isValid := isValidEmail(postData.Email)
		exists := emailExists(s.db.DB, postData.Email)
		if !isValid {
			http.Error(w, "Email not valid", http.StatusBadRequest)
			return
		} else if exists {
			http.Error(w, "Account with this email already exists", http.StatusBadRequest)
			return
		}
	}

	if postData.Password == "" {
		http.Error(w, "Password required", http.StatusBadRequest)
		return
	}

	if postData.DateOfBirth == "" {
		http.Error(w, "Date Of Birth required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(postData.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	hashPassword := string(hashedPassword)

	err = models.RegisterUser(s.db.DB, postData.Email, hashPassword, postData.FirstName, postData.LastName, postData.DateOfBirth, postData.Nickname, postData.AboutMe, postData.MimeType, postData.ImageData)
	if err != nil {
		fmt.Println(err)
	}

	// auto login after register

	valid, id, _ := models.ValidateLogin(s.db.DB, postData.Email, postData.Password)
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

func emailExists(db *sqlx.DB, email string) bool {
	tx, err := db.Beginx()
	if err != nil {
		return false
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	// Normalize the email address to lowercase to ensure case-insensitive comparison.
	email = strings.ToLower(email)

	// Prepare the SQL query.
	query := "SELECT 1 FROM users WHERE email = lower(?) LIMIT 1"

	// Execute the query and check if the email exists.
	var result int
	err = db.QueryRow(query, email).Scan(&result)

	if err == sql.ErrNoRows {
		// The email does not exist in the database.
		return false
	} else if err != nil {
		// An error occurred while executing the query.
		return false
	}

	// The email exists in the database.
	return true
}
