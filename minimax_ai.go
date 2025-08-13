package main

// MinimaxAI implementation
type MinimaxAI struct{}

func (ai MinimaxAI) GetMove(currentBoard [boardHeight][boardWidth]rune, currentPlayer rune) (int, int) {
	bestScore := -10000
	bestMove := [2]int{-1, -1}

	emptyCells := getEmptyCells(currentBoard)

	for _, cell := range emptyCells {
		boardCopy := currentBoard
		boardCopy[cell[1]][cell[0]] = currentPlayer

		score := minimax(boardCopy, 0, false, currentPlayer)

		if score > bestScore {
			bestScore = score
			bestMove = cell
		}
	}
	return bestMove[0], bestMove[1]
}

func (ai MinimaxAI) Name() string {
	return "Minimax AI"
}

func minimax(currentBoard [boardHeight][boardWidth]rune, depth int, isMaximizingPlayer bool, originalPlayer rune) int {
	opponent := 'X'
	if originalPlayer == 'X' {
		opponent = 'O'
	}

	if checkWin(currentBoard, originalPlayer) {
		return 10 - depth
	} else if checkWin(currentBoard, opponent) {
		return depth - 10
	} else if checkDraw(currentBoard) {
		return 0
	}

	emptyCells := getEmptyCells(currentBoard)

	if isMaximizingPlayer {
		bestScore := -10000
		for _, cell := range emptyCells {
			boardCopy := currentBoard
			boardCopy[cell[1]][cell[0]] = originalPlayer
			score := minimax(boardCopy, depth+1, false, originalPlayer)
			if score > bestScore {
				bestScore = score
			}
		}
		return bestScore
	} else {
		bestScore := 10000
		for _, cell := range emptyCells {
			boardCopy := currentBoard
			boardCopy[cell[1]][cell[0]] = opponent
			score := minimax(boardCopy, depth+1, true, originalPlayer)
			if score < bestScore {
				bestScore = score
			}
		}
		return bestScore
	}
}
