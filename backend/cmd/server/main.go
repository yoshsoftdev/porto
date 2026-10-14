package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/yoshsoftdev/porto/backend/internal/config"
	"github.com/yoshsoftdev/porto/backend/internal/database"
	"github.com/yoshsoftdev/porto/backend/internal/handlers"
	"github.com/yoshsoftdev/porto/backend/internal/seed"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	if cfg.AdminPassword == "" {
		cfg.AdminPassword = "changeme"
		log.Println("warning: ADMIN_PASSWORD is empty, using 'changeme'")
	}
	if cfg.JWTSecret == "" {
		cfg.JWTSecret = "dev-jwt-secret-change-me"
		log.Println("warning: JWT_SECRET is empty, using an insecure default")
	}

	db, err := database.Open(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := seed.Run(db, os.Getenv("PORTO_MD_PATH")); err != nil {
		log.Fatal("seed: ", err)
	}

	r := handlers.New(db, cfg)
	addr := ":" + cfg.Port
	log.Println("listening on", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
