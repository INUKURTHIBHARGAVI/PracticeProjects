package strategies

import (
	"math/rand"
	"tic-tac-toe/models"
	"time"
)

type clusterPlayingStrategy struct{}

func (c clusterPlayingStrategy) Play(board *models.Board) *models.Move {
	// Implement custom playing strategy
	rand.Seed(time.Now().UnixNano())
	availableCells := board.GetAvailableCells()
	randomIndex := rand.Intn(len(availableCells))
	randomCell := availableCells[randomIndex]
	return &models.Move{Row: randomCell.Row, Column: randomCell.Column}
}
