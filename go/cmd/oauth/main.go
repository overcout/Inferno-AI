package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/overcout/Inferno-AI/internal/config"
	"github.com/overcout/Inferno-AI/internal/oauth"
	"github.com/overcout/Inferno-AI/internal/store"
)

func main() {
	_ = godotenv.Load()
	cfg := config.LoadConfig()

	storeInstance, err := store.NewStore(cfg.DBDSN)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	oauth.InitOAuth(cfg)
	oauth.RegisterHandlers(storeInstance)

	log.Println("[OAUTH] Standalone server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
