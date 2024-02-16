package main

import (
	"fmt"
	"math/rand"
	"time"
)

type GameStatus string

const (
	IN_PROGRESS GameStatus = "IN_PROGRESS"
	DRAW        GameStatus = "DRAW"
	ENDED       GameStatus = "ENDED"
)

type GameSymbol string

const (
	X GameSymbol = "X"
	O GameSymbol = "O"
)

type User struct {
	Username string
	Email    string
	Photo    string
}

type Player interface {
	Play(board *Board) *Move
}

type Board struct {
	Cells [][]*BoardCell
}

type BoardCell struct {
	Row    int
	Column int
	Symbol GameSymbol
}

type Move struct {
	Row    int
	Column int
}

type HumanPlayer struct {
	Symbol GameSymbol
	User   *User
}

type Bot struct {
	Symbol          GameSymbol
	PlayingStrategy PlayingStrategy
}

type PlayingStrategy interface {
	Play(board *Board) *Move
}

type RandomPlayingStrategy struct{}

type clusterPlayingStrategy struct{}

type GameBuilder struct {
	Board         *Board
	Players       []Player
	Status        GameStatus
	NextPlayerIdx int
}

type Game struct {
	Board         *Board
	Players       []Player
	Status        GameStatus
	NextPlayerIdx int
}

type PlayerBuilder struct {
	Symbol GameSymbol
	Player Player
}

type UserBuilder struct {
	Username string
	email    string
	photo    []byte
}

func (h *HumanPlayer) Play(board *Board) *Move {
	var row, column int
	fmt.Println("Enter the row and column:")
	fmt.Scan(&row, &column)
	// row = 0
	// column = 0
	return &Move{Row: row, Column: column}
}

func (b *Bot) Play(board *Board) *Move {
	return b.PlayingStrategy.Play(board)
}

func (r *RandomPlayingStrategy) Play(board *Board) *Move {
	rand.Seed(time.Now().UnixNano())
	availableCells := board.GetAvailableCells()
	randomIndex := rand.Intn(len(availableCells))
	randomCell := availableCells[randomIndex]
	return &Move{Row: randomCell.Row, Column: randomCell.Column}
}

func (c *clusterPlayingStrategy) Play(board *Board) *Move {
	// Implement custom playing strategy
	return nil
}

func NewBoard(rows, columns int) *Board {
	cells := make([][]*BoardCell, rows)
	for i := 0; i < rows; i++ {
		cells[i] = make([]*BoardCell, columns)
		for j := 0; j < columns; j++ {
			cells[i][j] = &BoardCell{Row: i, Column: j}
		}
	}
	return &Board{Cells: cells}
}

func (b *Board) IsEmpty(row, column int) bool {
	return b.Cells[row][column].Symbol == ""
}

func (b *Board) GetAvailableCells() []*BoardCell {
	var availableCells []*BoardCell
	for _, row := range b.Cells {
		for _, cell := range row {
			if cell.Symbol == "" {
				availableCells = append(availableCells, cell)
			}
		}
	}
	return availableCells
}

func NewGameBuilder() *GameBuilder {
	return &GameBuilder{
		Status:        IN_PROGRESS,
		NextPlayerIdx: 0,
		Board:         &Board{},
		Players:       []Player{},
	}
}

func (gb *GameBuilder) WithDimensions(row int, col int) *GameBuilder {
	gb.Board = NewBoard(row, col)
	return gb
}

func (gb *GameBuilder) WithPlayer(p Player) *GameBuilder {
	gb.Players = append(gb.Players, p)
	return gb
}

func (gb *GameBuilder) Build() *Game {
	return &Game{
		Board:         gb.Board,
		Players:       gb.Players,
		Status:        gb.Status,
		NextPlayerIdx: gb.NextPlayerIdx,
	}
}

func NewPlayerBuilder() *PlayerBuilder {
	return &PlayerBuilder{}
}
func (pb *PlayerBuilder) WithSymbol(s GameSymbol) *PlayerBuilder {
	pb.Symbol = s
	return pb
}

func (pb *PlayerBuilder) WithUser(u *User) *PlayerBuilder {
	pb.Player = &HumanPlayer{Symbol: pb.Symbol, User: u}
	return pb
}

func (pb *PlayerBuilder) WithPlayingStrategy(ps PlayingStrategy) *PlayerBuilder {
	pb.Player = &Bot{Symbol: pb.Symbol, PlayingStrategy: ps}
	return pb
}

func (pb *PlayerBuilder) Build() Player {
	return pb.Player
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}
func (u *UserBuilder) WithUsername(username string) *UserBuilder {
	u.Username = username
	return u
}

func (u *UserBuilder) WithEmail(email string) *UserBuilder {
	u.email = email
	return u
}

func (u *UserBuilder) WithPhoto(photo []byte) *UserBuilder {
	u.photo = photo
	return u
}

func (u *UserBuilder) Build() *User {
	return &User{Username: u.Username, Email: u.email, Photo: string(u.photo)}
}

func (g *Game) GetNextPlayer() Player {
	return g.Players[g.NextPlayerIdx]
}

func (g *Game) MakeMove() {
	player := g.GetNextPlayer()
	move := player.Play(g.Board)

	if !g.Board.IsEmpty(move.Row, move.Column) {
		panic("Cell is not empty")
	}

	switch p := player.(type) {
	case *HumanPlayer:
		g.Board.Cells[move.Row][move.Column].Symbol = p.Symbol
	case *Bot:
		g.Board.Cells[move.Row][move.Column].Symbol = p.Symbol
	default:
		panic("Unknown player type")
	}

	if g.CheckWinner(g.Board.Cells[move.Row][move.Column].Symbol) {
		g.Status = ENDED
		fmt.Printf("win!\n")
		return
	}

	if g.CheckDraw() {
		g.Status = DRAW
		fmt.Printf("Game Drawn\n")
		return
	}

	g.NextPlayerIdx = (g.NextPlayerIdx + 1) % len(g.Players)
}

func (g *Game) CheckWinner(symbol GameSymbol) bool {
	return g.CheckRows(symbol) || g.CheckColumns(symbol) || g.CheckDiagonals(symbol)
}

func (g *Game) CheckRows(symbol GameSymbol) bool {
	for _, row := range g.Board.Cells {
		isWinner := true
		for _, cell := range row {
			if cell.Symbol != symbol {
				isWinner = false
				break
			}
		}
		if isWinner {
			return true
		}
	}
	return false
}

func (g *Game) CheckColumns(symbol GameSymbol) bool {
	for col := 0; col < len(g.Board.Cells[0]); col++ {
		win := true
		for row := 0; row < len(g.Board.Cells); row++ {
			if g.Board.Cells[row][col].Symbol != symbol {
				win = false
				break
			}
		}
		if win {
			return true
		}
	}
	return false
}

func (g *Game) CheckDiagonals(symbol GameSymbol) bool {
	win := true
	for i := 0; i < len(g.Board.Cells); i++ {
		if g.Board.Cells[i][i].Symbol != symbol {
			win = false
			break
		}
	}
	if win {
		return true
	}

	// Check top-right to bottom-left diagonal
	win = true
	for i := 0; i < len(g.Board.Cells); i++ {
		if g.Board.Cells[i][len(g.Board.Cells)-1-i].Symbol != symbol {
			win = false
			break
		}
	}
	return win
}

func (g *Game) CheckDraw() bool {
	return len(g.Board.GetAvailableCells()) == 0
}

func (g *Game) PrintBoard() {
	for _, row := range g.Board.Cells {
		for _, cell := range row {
			if cell.Symbol == "" {
				fmt.Print("| - ")
			} else {
				fmt.Printf("| %s ", cell.Symbol)
			}
		}
		fmt.Println("|")
	}
}

func (g *Game) WithPlayer(p Player) {
	g.Players = append(g.Players, p)
}

func main() {
	game := createGame()
	fmt.Println("Game Started")

	for game.Status == IN_PROGRESS {
		nextPlayer := game.GetNextPlayer()

		switch p := nextPlayer.(type) {
		case *HumanPlayer:
			fmt.Printf("Next Move is of Player %s\n", p.Symbol)
			game.MakeMove()
		case *Bot:
			fmt.Printf("Bot %s is making its move...\n", p.Symbol)
			game.MakeMove()

		default:
			panic("Unknown player type")
		}

		// Check status again before printing the board to ensure accurate state
		if game.Status == IN_PROGRESS {
			game.PrintBoard()
		}
	}

}

func createGame() *Game {
	user := NewUserBuilder().WithUsername("Bhargavi").WithEmail("hey@gmail.com").WithPhoto([]byte("photo")).Build()
	human := NewPlayerBuilder().WithSymbol("0").WithUser(user).Build()
	bot := NewPlayerBuilder().WithSymbol("X").WithPlayingStrategy(&RandomPlayingStrategy{}).Build()
	game := NewGameBuilder().WithDimensions(3, 3).WithPlayer(human).WithPlayer(bot).Build()
	return game
}
