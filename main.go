package main

import (
	"github.com/jbdoumenjou/mygoserver/internal/db"
	"log"
)

func main() {
	db, err := db.NewDB("database.json")
	if err != nil {
		panic(err)
	}

	router := NewRouter(db)
	log.Fatal(NewWebServer(":8080", router).Start())
}
