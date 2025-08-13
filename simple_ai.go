package main

import (
	"math/rand"
)

// SimpleAI implementation
type SimpleAI struct{}

func (ai SimpleAI) GetMove(currentBoard [boardHeight][boardWidth]rune, currentPlayer rune) (int, int) {
	// 1. Check for a winning move
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			if currentBoard[y][x] == ' ' {
				currentBoard[y][x] = currentPlayer
				if checkWin(currentBoard, currentPlayer) {
					currentBoard[y][x] = ' ' // backtrack
					return x, y
				}
				currentBoard[y][x] = ' ' // backtrack
			}
		}
	}

	// 2. Block opponent's winning move
	opponent := 'X'
	if currentPlayer == 'X' {
		opponent = 'O'
	}

	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			if currentBoard[y][x] == ' ' {
				currentBoard[y][x] = opponent
				if checkWin(currentBoard, opponent) {
					currentBoard[y][x] = ' ' // backtrack
					return x, y
				}
				currentBoard[y][x] = ' ' // backtrack
			}
		}
	}

	// 3. Pick a random empty spot
	emptyCells := getEmptyCells(currentBoard)

	if len(emptyCells) > 0 {
		move := emptyCells[rand.Intn(len(emptyCells))]
		return move[0], move[1]
	}
	return -1, -1 // Should not happen in a valid game
}

func (ai SimpleAI) Name() string {
	return "Simple AI"
}
