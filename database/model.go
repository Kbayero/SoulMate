package database

import (
	"time"
)

type BaseModel struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Result struct {
	BaseModel
	Player1Name   string `json:"player1_name"`
	Player2Name   string `json:"player2_name"`
	Question      string `json:"question"`
	Player1Answer string `json:"player1_answer"`
	Player2Answer string `json:"player2_answer"`
	IsCorrect     bool   `json:"is_correct"`
}

func Migrate() {
	db := GetDB()
	err := db.Migrate(&Result{})
	if err != nil {
		panic(err)
	}
}
