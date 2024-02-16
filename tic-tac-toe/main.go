package main

import (
	"fmt"
	"tic-tac-toe/models"
	"tic-tac-toe/strategies"
)

func main() {
	game := createGame()
	fmt.Println("Game Started")

	for game.Status == models.IN_PROGRESS {
		nextPlayer := game.GetNextPlayer()

		switch p := nextPlayer.(type) {
		case *models.HumanPlayer:
			fmt.Printf("Next Move is of Player %s\n", p.Symbol)
			game.MakeMove()
		case *models.Bot:
			fmt.Printf("Bot %s is making its move...\n", p.Symbol)
			game.MakeMove()

		default:
			panic("Unknown player type")
		}

		// Check status again before printing the board to ensure accurate state
		if game.Status == models.IN_PROGRESS {
			game.PrintBoard()
		}
	}

}

func createGame() *models.Game {
	user := models.NewUserBuilder().WithUsername("Bhargavi").WithEmail("hey@gmail.com").WithPhoto([]byte("photo")).Build()
	human := models.NewPlayerBuilder().WithSymbol("0").WithUser(user).Build()
	bot := models.NewPlayerBuilder().WithSymbol("X").WithPlayingStrategy(&strategies.RandomPlayingStrategy{}).Build()
	game := models.NewGameBuilder().WithDimensions(3, 3).WithPlayer(human).WithPlayer(bot).Build()
	return game
}
