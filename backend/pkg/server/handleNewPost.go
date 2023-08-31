package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"social-network/pkg/models"
	mw "social-network/pkg/server/middleware"
	"strconv"
	"strings"
)

func (s *Server) newPostHandler(w http.ResponseWriter, r *http.Request) {
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
	postData.Title = r.FormValue("title")
	postData.PostTarget = r.FormValue("post_target")
	postData.Content = r.FormValue("content")
	privacy := r.FormValue("privacy")
	privacyInt, _ := strconv.Atoi(privacy)
	groupid := r.FormValue("groupid")
	groupidInt, _ := strconv.Atoi(groupid)
	postData.GroupID = groupidInt
	postData.Privacy = privacyInt
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

	followersJSON := r.FormValue("followers")
	if len(followersJSON) > 0 {
		decoder := json.NewDecoder(strings.NewReader(followersJSON))
		err := decoder.Decode(&postData.Followers)
		if err != nil {
			http.Error(w, "Unable to parse followers data", http.StatusInternalServerError)
			return
		}
	}

	log.Printf("Received post data:\n")
	log.Printf("Title: %s\n", postData.Title)
	log.Printf("Content: %s\n", postData.Content)
	log.Printf("Target: %s\n", postData.PostTarget)
	log.Printf("Privacy: %d\n", postData.Privacy)
	log.Printf("Created By: %d\n", postData.CreatedBy)
	log.Printf("Followers: %d\n", postData.CreatedBy)
	log.Printf("Group id: %d\n", postData.GroupID)
	if postData.Privacy == 2 && len(postData.Followers) > 0 {
		log.Println("Received post followers:")
		log.Println(postData.Followers)
	}

	result, err := models.NewPost(postData, s.db.DB)
	if err != nil {
		log.Printf("error adding new post %v", err)
	}

	var url string
	if postData.PostTarget == "regular_post" {
		url = "/post/" + strconv.Itoa(result.PostID)
	} else if postData.PostTarget == "group_post" {
		url = "/groups/post/" + strconv.Itoa(result.PostID)
	}
	response := struct {
		PostID  int    `json:"post_id"`
		Message string `json:"message"`
		Url     string `json:"url"`
	}{
		PostID:  result.PostID,
		Message: "Post created successfully",
		Url:     url,
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
