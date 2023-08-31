package main

import (
	"log"
	"social-network/pkg/server"
	"social-network/pkg/utils"
)

func main() {
	srv := server.New()
	port := utils.GetEnv("BACKEND_ADDR")

	log.Printf("Starting server on http://localhost%s", port)
	if err := srv.Run(); err != nil {
		log.Fatal("Error running server:", err)
	}
}
