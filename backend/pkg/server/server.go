package server

import (
	"log"
	"net/http"
	"social-network/pkg/db"
	"social-network/pkg/models"
	"social-network/pkg/utils"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	db        *db.Database
	secretKey string
	addr      string

	connsMu sync.RWMutex
	conns   map[int]*websocket.Conn
	writeMu sync.Mutex

	broadcast chan models.Message
}

func New() *Server {
	db, err := db.Open()
	if err != nil {
		log.Fatalf("can't open db: %v", err)
	}

	port := utils.GetEnv("BACKEND_ADDR")
	s := &Server{
		db:   db,
		addr: port,

		connsMu: sync.RWMutex{},
		conns:   make(map[int]*websocket.Conn),

		writeMu: sync.Mutex{},

		broadcast: make(chan models.Message, 10),
	}

	return s
}

func (s *Server) Run() error {
	defer s.db.Close()
	server := &http.Server{
		Handler:      s.router(),
		ReadTimeout:  time.Duration(5) * time.Second,
		WriteTimeout: time.Duration(10) * time.Second,
	}

	err := s.db.RunMigrations()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Seeder(s.db.DB)

	go s.ListenForMessages(s.broadcast)
	err = server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) AddConn(userID int, conn *websocket.Conn) {
	s.connsMu.Lock()
	defer s.connsMu.Unlock()

	s.conns[userID] = conn
}

func (s *Server) RemoveConn(userID int) {
	s.connsMu.Lock()
	defer s.connsMu.Unlock()

	delete(s.conns, userID)
}

func (s *Server) GetConn(userID int) *websocket.Conn {
	s.connsMu.Lock()
	defer s.connsMu.Unlock()

	return s.conns[userID]
}
