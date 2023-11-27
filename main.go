package main

import (
	"log"
	"os"

	"github.com/jbdoumenjou/mygoserver/internal/api/token"

	"github.com/jbdoumenjou/mygoserver/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	apiKey := os.Getenv("API_KEY")

	db, err := db.NewDB("database.json")
	if err != nil {
		panic(err)
	}

	tokenManager := token.NewManager(jwtSecret, apiKey)
	router := NewRouter(db, tokenManager)
	log.Fatal(NewWebServer(":8080", router).Start())
}
