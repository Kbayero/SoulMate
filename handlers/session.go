package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/KbaYero/SoulMate/services"
	"github.com/google/uuid"
	"golang.org/x/exp/rand"
)

func SubmitSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/start", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	sessionID := r.FormValue("session_id")

	if name == "" {
		http.Redirect(w, r, "/start", http.StatusSeeOther)
		return
	}

	services.Mutex.Lock()
	defer services.Mutex.Unlock()

	if sessionID == "" {
		newSessionID := generateSessionID()
		player1ID := uuid.New().String()

		player1 := services.Player{
			ID:              player1ID,
			Name:            name,
			Answers:         []string{},
			CurrentQuestion: 1,
		}

		game := &services.Game{
			ID:      newSessionID,
			Player1: player1,
			Status:  "waiting",
		}

		services.Games[newSessionID] = game

		log.Printf("Nueva sesión creada: %s por Player1: %s\n", newSessionID, name)

		templates.ExecuteTemplate(w, "waiting.html", map[string]interface{}{
			"SessionID":       newSessionID,
			"PlayerID":        player1ID,
			"ShowSessionInfo": true,
		})
		return
	} else {
		game, exists := services.Games[sessionID]
		if !exists {
			log.Printf("Intento de unión a sesión inexistente: %s\n", sessionID)
			templates.ExecuteTemplate(w, "error.html", map[string]interface{}{
				"Message": "La sesión proporcionada no existe.",
			})
			return
		}

		if game.Status != "waiting" {
			log.Printf("Intento de unión a sesión no disponible: %s\n", sessionID)
			templates.ExecuteTemplate(w, "error.html", map[string]interface{}{
				"Message": "La sesión ya está en progreso o ha finalizado.",
			})
			return
		}

		// Crear Player2
		player2ID := uuid.New().String()
		player2 := services.Player{
			ID:              player2ID,
			Name:            name,
			Answers:         []string{},
			CurrentQuestion: 1,
		}

		game.Player2 = player2
		game.Status = "in_progress"

		http.Redirect(w, r, fmt.Sprintf("/game?session_id=%s&player_id=%s", sessionID, player2ID), http.StatusSeeOther)
		return
	}
}

func generateSessionID() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	const length = 8

	rand.Seed(uint64(time.Now().UnixNano()))

	randomString := func() string {
		b := make([]byte, length)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}
		return string(b)
	}

	for {
		sessionID := randomString()
		log.Printf("Generando sessionID: %s\n", sessionID)

		_, exists := services.Games[sessionID]

		if !exists {
			log.Printf("sessionID %s es único y se utilizará.\n", sessionID)
			return sessionID
		}
		log.Printf("SessionID %s ya existe, generando nuevo...\n", sessionID)
	}
}
