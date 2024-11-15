package handlers

import (
	"log"
	"net/http"

	"github.com/KbaYero/SoulMate/services"
)

func ResultHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	playerID := r.URL.Query().Get("player_id")

	if sessionID == "" || playerID == "" {
		log.Println("Error: session_id o player_id faltante en resultHandler")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	services.Mutex.Lock()
	game, exists := services.Games[sessionID]
	services.Mutex.Unlock()

	if !exists {
		log.Printf("Error: Juego no encontrado en resultHandler para session_id %s\n", sessionID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if game.Status != "finished" {
		log.Printf("Error: Juego no está finalizado en resultHandler para session_id %s\n", sessionID)
		http.Redirect(w, r, "/game?session_id="+sessionID+"&player_id="+playerID, http.StatusSeeOther)
		return
	}

	questions := []string{}
	for _, question := range game.Questions {
		questions = append(questions, question.Text)
	}

	if len(game.Results) == 0 {
		game.Results, game.Percentage = services.GetResponses(questions, game.Player1.Answers, game.Player2.Answers, game.Player1.Name, game.Player2.Name)
	}

	game.ResultsShown++
	log.Printf("Player %s ha visto los resultados. Total Players que han visto: %d\n", playerID, game.ResultsShown)

	if game.ResultsShown >= 2 {
		delete(services.Games, sessionID)
		log.Printf("Juego %s eliminado después de que ambos jugadores vieron los resultados.\n", sessionID)
	}

	templates.ExecuteTemplate(w, "result.html", map[string]interface{}{
		"Percentage":    game.Percentage,
		"Player1":       game.Player1.Name,
		"Player2":       game.Player2.Name,
		"ResultDetails": game.Results,
	})
}
