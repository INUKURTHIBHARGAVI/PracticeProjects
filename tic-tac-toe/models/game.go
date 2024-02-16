package models

import "fmt"

type GameSymbol string

const (
	X GameSymbol = "X"
	O GameSymbol = "O"
)

type GameStatus string

const (
	IN_PROGRESS GameStatus = "IN_PROGRESS"
	DRAW        GameStatus = "DRAW"
	ENDED       GameStatus = "ENDED"
)

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
