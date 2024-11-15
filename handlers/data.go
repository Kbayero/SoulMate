package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/KbaYero/SoulMate/database"
	"github.com/KbaYero/SoulMate/services"
)

func RingHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("auth")
	expectedAuth := os.Getenv("AUTH_KEY")

	if authHeader != expectedAuth {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	ring := r.FormValue("ring_question")
	if ring == "yes" {
		services.AskForRing = true
	} else {
		services.AskForRing = false
	}
	w.WriteHeader(http.StatusOK)
}

func GetDataHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("auth")
	expectedAuth := os.Getenv("AUTH_KEY")

	if authHeader != expectedAuth {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var results []database.Result
	if err := database.GetDB().GetAll(&results); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
