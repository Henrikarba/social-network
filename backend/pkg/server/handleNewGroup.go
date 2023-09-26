package server

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/pkg/models"
	mw "social-network/pkg/server/middleware"
)

func (s *Server) newGroup(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(mw.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // set maxFormSize appropriately
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	var postData models.Group
	postData.CreatorID = id
	postData.Title = r.FormValue("title")
	postData.Description = r.FormValue("description")

	gid, err := models.NewGroup(s.db.DB, postData)
	if err != nil {
		log.Printf("new group: %v", err)
	}

	response := struct {
		Message string `json:"message"`
		GroupID int    `json:"group_id"`
	}{

		Message: "Group created successfully",
		GroupID: gid,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	conn, ok := s.conns[id]
	if ok {
		data, _ := models.GetAuthenticatedUserDate(s.db.DB, id)
		s.writeMu.Lock()
		conn.WriteJSON(data)
		s.writeMu.Unlock()
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
