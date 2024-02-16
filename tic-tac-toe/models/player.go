package models

import (
	"fmt"
)

type Player interface {
	Play(board *Board) *Move
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

type PlayerBuilder struct {
	Symbol GameSymbol
	Player Player
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
