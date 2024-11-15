// main.go
package main

import (
	"log"
	"net/http"

	"github.com/KbaYero/SoulMate/database"
	"github.com/KbaYero/SoulMate/handlers"
)

func main() {
	database.Migrate()
	handlers.InitTemplates()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/start", handlers.StartHandler)
	http.HandleFunc("/submit_session", handlers.SubmitSessionHandler)
	http.HandleFunc("/game", handlers.GameHandler)
	http.HandleFunc("/result", handlers.ResultHandler)
	http.HandleFunc("/status", handlers.StatusHandler)
	http.HandleFunc("/error", handlers.ErrorHandler)
	http.HandleFunc("/ask-for-ring", handlers.RingHandler)
	http.HandleFunc("/get-data", handlers.GetDataHandler)

	log.Println("Server started: http://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
