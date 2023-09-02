package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"social-network/pkg/models"
	mw "social-network/pkg/server/middleware"
	"strconv"
)

func (s *Server) newComment(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(mw.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	var postData models.PostRequest
	postData.PostTarget = r.FormValue("post_target")
	postData.Content = r.FormValue("content")
	groupid := r.FormValue("group_id")
	groupidInt, _ := strconv.Atoi(groupid)
	postData.GroupID = groupidInt
	postid := r.FormValue("post_id")
	postidInt, _ := strconv.Atoi(postid)
	postData.ID = postidInt
	postData.CreatedBy = id
	file, header, err := r.FormFile("image")
	if err == nil {
		postData.MimeType = header.Header.Get("Content-Type")
		postData.ImageData, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read file data", http.StatusInternalServerError)
			return
		}
	}

	if len(postData.ImageData) > 0 {
		log.Printf("Received image with MIME type: %s\n", postData.MimeType)
	}

	comment, err := models.NewComment(postData, s.db.DB)
	if err != nil {
		log.Printf("%v", err)
	}
	response := struct {
		Message string         `json:"message"`
		Comment models.Comment `json:"comment"`
	}{

		Message: "Post created successfully",
		Comment: *comment,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
