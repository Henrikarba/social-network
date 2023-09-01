package server

import (
	"net/http"
	mw "social-network/pkg/server/middleware"
)

func (s *Server) router() *http.ServeMux {
	r := http.NewServeMux()

	// Handle the login route with CORS middleware
	r.Handle("/login", mw.Cors(http.HandlerFunc(s.handleLogin)))

	// New posts, logins, registration, etc.
	authMiddleware := mw.AuthenticatedUser(s.db.DB)
	r.Handle("/new/post", authMiddleware(http.HandlerFunc(s.newPostHandler)))

	// File server
	fileHandler := http.StripPrefix("/images/", http.FileServer(http.Dir("./pkg/db/files/images")))
	r.Handle("/images/", authMiddleware(fileHandler))

	// WebSocket
	r.Handle("/ws", authMiddleware(http.HandlerFunc(s.wsHandler)))

	return r
}