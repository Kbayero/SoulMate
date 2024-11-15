package services

import (
	"sync"
	"time"

	"github.com/KbaYero/SoulMate/database"
	"golang.org/x/exp/rand"
)

type Player struct {
	ID              string
	Name            string
	Answers         []string
	CurrentQuestion int
}

type Game struct {
	ID           string
	Player1      Player
	Player2      Player
	Questions    []Question
	Results      []database.Result
	Percentage   int
	Status       string // "round1", "round2", "waiting", "in_progress", "finished"
	ResultsShown int
}

type StatusResponse struct {
	Status string `json:"status"`
}

type Question struct {
	Text  string
	Photo string
}

var (
	Games          = make(map[string]*Game)
	Mutex          = &sync.Mutex{}
	TotalQuestions = 10
	AskForRing     = false
)

func GetQuestions() ([]Question, error) {
	allQuestions := []Question{
		{"Where was he/she born?", "birthplace.jpg"},
		{"What is his/her favorite memory with you?", "memory.jpg"},
		{"What are his/her parents' names?", "parents.jpg"},
		{"What is his/her most important family tradition?", "family.jpg"},
		{"Did he/she have any pets during childhood?", "pets.jpg"},
		{"What was the name of the first company where he/she worked?", "first_job.jpg"},
		{"Has he/she lived in any other country besides where he/she grew up?", "other_city.jpg"},
		{"What is his/her favorite food?", "food.jpg"},
		{"Does he/she prefer coffee or tea?", "coffee_tea.jpg"},
		{"What is his/her favorite movie?", "movie.jpg"},
		{"Does he/she prefer the beach or the mountains?", "beach_mountains.jpg"},
		{"What is his/her favorite season of the year?", "season.jpg"},
		{"What hobby or pastime does he/she enjoy the most?", "hobby.jpg"},
		{"What goal would he/she like to achieve the next 5 year?", "dream.jpg"},
		{"Is there any place in the world he/she would love to visit?", "place.jpg"},
		{"What is his/her love language?", "love_language.jpg"},
		{"How important is family in his/her life?", "family_importance.jpg"},
		{"What adventure would he/she like to experience in the future?", "adventure.jpg"},
		{"Would he/she prefer to live in a house or apartment?", "house_apartment.jpg"},
		{"What is his/her favorite way to show affection?", "affection.jpg"},
		{"Does he/she have any special talents?", "talents.jpg"},
		{"What languages does he/she speak?", "languages.jpg"},
		{"What is his/her favorite TV series?", "tv_series.jpg"},
		{"What does he/she like to do in his/her free time?", "free_time.jpg"},
		{"Does he/she have any idols or people he/she admires?", "idols.jpg"},
		{"Does he/she prefer living in the city or in the countryside?", "city_countryside.jpg"},
		{"What is his/her favorite sport to play?", "sport.jpg"},
		{"Does he/she have any fears or phobias?", "fears.jpg"},
		{"Does he/she prefer reading or watching movies?", "reading_movies.jpg"},
		{"What is his/her preferred way to travel?", "travel.jpg"},
		{"Does he/she prefer dogs or cats?", "dogs_cats.jpg"},
		{"Does he/she prefer going out with friends or spending time alone?", "friends_alone.jpg"},
		{"What is his/her preferred form of exercise?", "exercise.jpg"},
		{"Does he/she prefer the cinema or the theater?", "cinema_theater.jpg"},
		{"How does he/she usually start his/her day?", "morning.jpg"},
		{"Does he/she have any health-related goals?", "health.jpg"},
		{"Does he/she prefer silence or music?", "silence_music.jpg"},
	}

	rand.Seed(uint64(time.Now().UnixNano()))
	rand.Shuffle(len(allQuestions), func(i, j int) {
		allQuestions[i], allQuestions[j] = allQuestions[j], allQuestions[i]
	})

	selectedQuestions := allQuestions[:TotalQuestions]

	if AskForRing {
		rand.Seed(uint64(time.Now().UnixNano()))
		indexToRemove := rand.Intn(len(selectedQuestions))
		selectedQuestions = append(selectedQuestions[:indexToRemove], selectedQuestions[indexToRemove+1:]...)

		selectedQuestions = append(selectedQuestions, Question{
			Text:  "What size wedding ring does he/she wear?",
			Photo: "ring_size.jpg",
		})

		rand.Seed(uint64(time.Now().UnixNano()))
		rand.Shuffle(len(selectedQuestions), func(i, j int) {
			selectedQuestions[i], selectedQuestions[j] = selectedQuestions[j], selectedQuestions[i]
		})
	}

	return selectedQuestions, nil
}
