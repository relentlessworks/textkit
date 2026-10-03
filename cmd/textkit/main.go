package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/relentlessworks/textkit/internal/api"
	"github.com/relentlessworks/textkit/internal/config"
)

func main() {
	cfg := config.Load()
	handler := api.New(cfg.NoAuth)
	mux := handler.Routes()

	log.Printf("textkit starting on %s (no-auth=%v)", cfg.Addr, cfg.NoAuth)
	fmt.Printf("textkit — text analysis and manipulation service\n")
	fmt.Printf("Listening on %s\n", cfg.Addr)
	fmt.Printf("GET /help for usage\n")

	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
