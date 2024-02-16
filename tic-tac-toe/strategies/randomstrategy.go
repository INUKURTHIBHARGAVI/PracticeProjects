package strategies

import (
	"math/rand"
	"tic-tac-toe/models"
	"time"
)

type RandomPlayingStrategy struct{}

func (r RandomPlayingStrategy) Play(board *models.Board) *models.Move {
	rand.Seed(time.Now().UnixNano())
	availableCells := board.GetAvailableCells()
	randomIndex := rand.Intn(len(availableCells))
	randomCell := availableCells[randomIndex]
	return &models.Move{Row: randomCell.Row, Column: randomCell.Column}
}
