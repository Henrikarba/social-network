package server

import (
	"net/http"
	mw "social-network/pkg/server/middleware"

	"github.com/gorilla/mux"
)

func (s *Server) router() *mux.Router {
	sm := mux.NewRouter()
	sm.Use(mw.Cors)
	authR := sm.Methods(http.MethodPost).Subrouter()
	authR.HandleFunc("/login", s.handleLogin)

	// New posts, logins, registration etc
	authMiddleware := mw.AuthenticatedUser(s.db.DB)
	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/new/post", s.newPostHandler)
	postR.Use(authMiddleware)

	// File server
	fileR := sm.Methods(http.MethodGet).Subrouter()
	fileHandler := http.StripPrefix("/images/", http.FileServer(http.Dir("./pkg/db/files/images")))
	fileR.PathPrefix("/images/").Handler(fileHandler)
	fileR.Use(authMiddleware)

	// WebSocket
	wsR := sm.Methods(http.MethodGet).Subrouter()
	wsR.HandleFunc("/ws", s.wsHandler)
	wsR.Use(authMiddleware)

	return sm
}
