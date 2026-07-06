package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Andrew1996-la/backend-diplom/pkg/api"
	"github.com/Andrew1996-la/backend-diplom/pkg/db"
)

func main() {
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		dbPath = "scheduler.db"
	}
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	api.Init()

	webDir := "web"

	if err := db.Init(dbPath); err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("server started on :%s", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
