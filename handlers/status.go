package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/KbaYero/SoulMate/services"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")

	if sessionID == "" {
		http.Error(w, "ID de juego faltante", http.StatusBadRequest)
		return
	}

	services.Mutex.Lock()
	game, exists := services.Games[sessionID]
	services.Mutex.Unlock()

	var response services.StatusResponse
	if !exists {
		response = services.StatusResponse{
			Status: "finished",
		}
	} else {
		response = services.StatusResponse{
			Status: game.Status,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
