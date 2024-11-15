package handlers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/KbaYero/SoulMate/services"
)

func GameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sessionID := r.URL.Query().Get("session_id")
	playerID := r.URL.Query().Get("player_id")

	if sessionID == "" || playerID == "" {
		log.Println("Error: missing session_id or player_id")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	services.Mutex.Lock()
	game, exists := services.Games[sessionID]
	services.Mutex.Unlock()

	if !exists {
		log.Printf("Error: game not found for session_id %s\n", sessionID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if len(game.Questions) == 0 {
		questions, err := services.GetQuestions()
		if err != nil {
			log.Println("Error getting questions:", err)
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}
		services.Mutex.Lock()
		game.Questions = questions
		services.Mutex.Unlock()
	}

	var currentPlayer *services.Player

	if game.Player1.ID == playerID {
		currentPlayer = &game.Player1
	} else if game.Player2.ID == playerID {
		currentPlayer = &game.Player2
	} else {
		log.Printf("Error: PlayerID %s is not from session %s\n", playerID, sessionID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	questions := game.Questions

	if r.Method == http.MethodPost {
		answer := r.FormValue("answer")
		if answer != "" {
			safeAnswer := template.HTMLEscapeString(answer)

			services.Mutex.Lock()
			currentPlayer.Answers = append(currentPlayer.Answers, answer)
			currentPlayer.CurrentQuestion++
			services.Mutex.Unlock()

			log.Printf("Player %s answered the question %d with: %s\n", playerID, currentPlayer.CurrentQuestion-1, safeAnswer)

			if currentPlayer.CurrentQuestion > len(questions) {
				services.Mutex.Lock()
				if game.Player1.CurrentQuestion > len(questions) && game.Player2.CurrentQuestion > len(questions) {
					game.Status = "finished"
					log.Printf("Game %s finished.\n", sessionID)
				}
				services.Mutex.Unlock()

				if game.Status == "finished" {
					http.Redirect(w, r, "/result?session_id="+sessionID+"&player_id="+playerID, http.StatusSeeOther)
					return
				} else {
					templates.ExecuteTemplate(w, "waiting.html", map[string]interface{}{
						"SessionID":       sessionID,
						"PlayerID":        playerID,
						"ShowSessionInfo": false,
					})
					return
				}
			}
		}
	}

	if currentPlayer.CurrentQuestion > len(questions) {
		if game.Status == "finished" {
			http.Redirect(w, r, "/result?session_id="+sessionID+"&player_id="+playerID, http.StatusSeeOther)
			return
		} else {
			templates.ExecuteTemplate(w, "waiting.html", map[string]interface{}{
				"SessionID":       sessionID,
				"PlayerID":        playerID,
				"ShowSessionInfo": false,
			})
			return
		}
	}

	if currentPlayer.CurrentQuestion-1 >= len(questions) {
		currentPlayer.CurrentQuestion = len(questions)
	}

	question := questions[currentPlayer.CurrentQuestion-1]

	subtitle := ""
	if game.Player1.ID == playerID {
		subtitle = fmt.Sprintf("Answer the following questions to compare them with the answers of %s.", game.Player2.Name)
	} else {
		subtitle = fmt.Sprintf("How much do you know %s?", game.Player1.Name)
	}

	templates.ExecuteTemplate(w, "game.html", map[string]interface{}{
		"SessionID":       sessionID,
		"PlayerID":        playerID,
		"Player":          getPlayerRole(game, playerID),
		"Question":        question.Text,
		"Photo":           question.Photo,
		"Subtitle":        subtitle,
		"CurrentQuestion": currentPlayer.CurrentQuestion,
	})
}

func getPlayerRole(game *services.Game, playerID string) string {
	if game.Player1.ID == playerID {
		return "Player1"
	}
	return "Player2"
}
