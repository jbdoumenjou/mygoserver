package main

import (
	"log"
)

func main() {
	router := NewRouter()
	log.Fatal(NewWebServer(":8080", router).Start())
}
