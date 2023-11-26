package main

import (
	"log"
	"os"

	"github.com/jbdoumenjou/mygoserver/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	jwtSecret := os.Getenv("JWT_SECRET")

	db, err := db.NewDB("database.json")
	if err != nil {
		panic(err)
	}

	router := NewRouter(db, ApiConfig{JWTSecret: jwtSecret})
	log.Fatal(NewWebServer(":8080", router).Start())
}
