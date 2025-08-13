package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	boardWidth  = 3
	boardHeight = 3
	pvp         = 0 // Player vs Player
	pvc         = 1 // Player vs Computer
	cvc         = 2 // Computer vs Computer
)

// ComputerPlayer interface for different AI implementations
type ComputerPlayer interface {
	GetMove(board [boardHeight][boardWidth]rune, player rune) (int, int)
	Name() string
}

// Player interface for both human and AI players
type Player interface {
	GetMove(currentBoard [boardHeight][boardWidth]rune, currentPlayer rune) (int, int, error)
	Name() string
	IsHuman() bool
}

var (
	board        [boardHeight][boardWidth]rune
	cursorX      = 0
	cursorY      = 0
	player       = 'X'
	gameOver     = false
	winner       = ' '
	message      = ""
	gameMode     = -1 // -1: menu, 0: pvp, 1: pvc, 2: cvc

	playerX Player
	playerO Player

	menuSelector = 0
)

// HumanPlayer implementation
type HumanPlayer struct{}

func (hp HumanPlayer) GetMove(currentBoard [boardHeight][boardWidth]rune, currentPlayer rune) (int, int, error) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				cursorY = (cursorY - 1 + boardHeight) % boardHeight
			case termbox.KeyArrowDown:
				cursorY = (cursorY + 1) % boardHeight
			case termbox.KeyArrowLeft:
				cursorX = (cursorX - 1 + boardWidth) % boardWidth
			case termbox.KeyArrowRight:
				cursorX = (cursorX + 1) % boardWidth
			case termbox.KeyEnter:
				if currentBoard[cursorY][cursorX] == ' ' {
					return cursorX, cursorY, nil
				}
			case termbox.KeyEsc, termbox.KeyCtrlC:
				gameMode = -1 // Go back to menu
				resetGame()
				return -1, -1, fmt.Errorf("user exited to menu")
			}
			draw() // Redraw after cursor movement
		case termbox.EventError:
			return -1, -1, ev.Err
		}
	}
}

func (hp HumanPlayer) Name() string {
	return "Human"
}

func (hp HumanPlayer) IsHuman() bool {
	return true
}

// ComputerPlayerWrapper wraps a ComputerPlayer to implement the Player interface
type ComputerPlayerWrapper struct {
	AI ComputerPlayer
}

func (cpw ComputerPlayerWrapper) GetMove(currentBoard [boardHeight][boardWidth]rune, currentPlayer rune) (int, int, error) {
	time.Sleep(500 * time.Millisecond) // Delay for AI moves
	x, y := cpw.AI.GetMove(currentBoard, currentPlayer)
	return x, y, nil
}

func (cpw ComputerPlayerWrapper) Name() string {
	return cpw.AI.Name()
}

func (cpw ComputerPlayerWrapper) IsHuman() bool {
	return false
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	rand.Seed(time.Now().UnixNano())
	termbox.SetInputMode(termbox.InputEsc)
	resetGame()

	for {
		if gameMode == -1 {
			drawMenu()
			handleMenuEvent()
		} else {
			draw()
			if gameOver {
				switch ev := termbox.PollEvent(); ev.Type {
				case termbox.EventKey:
					if ev.Key == termbox.KeyEnter || ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
						gameMode = -1 // Go back to menu
						resetGame()
					}
				case termbox.EventError:
					panic(ev.Err)
				}
				continue
			}

			var x, y int
			var err error

			// Get move from current player
			if player == 'X' {
				x, y, err = playerX.GetMove(board, player)
			} else { // player == 'O'
				x, y, err = playerO.GetMove(board, player)
			}

			if err != nil {
				if err.Error() == "user exited to menu" {
					continue // Go back to menu loop
				}
				panic(err)
			}

			// Apply move if valid
			if board[y][x] == ' ' {
				board[y][x] = player
				endTurn()
			}
			draw()
		}
	}
}

func handleMenuEvent() {
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyArrowUp:
			menuSelector = (menuSelector - 1 + 3) % 3
		case termbox.KeyArrowDown:
			menuSelector = (menuSelector + 1) % 3
		case termbox.KeyEnter:
			gameMode = menuSelector
			resetGame()
		case termbox.KeyEsc, termbox.KeyCtrlC:
			termbox.Close()
			panic("exiting") // A bit abrupt, but effective
		}
	case termbox.EventError:
		panic(ev.Err)
	}
}



func endTurn() {
	if checkWin(board, player) {
		winner = player
		gameOver = true
		if gameMode == pvc && player == 'O' {
			message = "Computer wins! Press Enter to continue."
		} else if gameMode == cvc {
			if winner == 'X' {
				message = "SimpleAI wins! Press Enter to continue."
			} else if winner == 'O' {
				message = "MinimaxAI wins! Press Enter to continue."
			}
		} else {
			message = fmt.Sprintf("Player %c wins! Press Enter to continue.", player)
		}
	} else if checkDraw(board) {
		gameOver = true
		message = "It's a draw! Press Enter to continue."
	} else {
		if player == 'X' {
			player = 'O'
		} else {
			player = 'X'
		}
	}
}

func resetGame() {
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			board[y][x] = ' '
		}
	}
	cursorX = 0
	cursorY = 0
	player = 'X'
	gameOver = false
	winner = ' '
	message = ""

	// Initialize players based on game mode
	switch gameMode {
	case pvp:
		playerX = HumanPlayer{}
		playerO = HumanPlayer{}
	case pvc:
		playerX = HumanPlayer{}
		playerO = ComputerPlayerWrapper{AI: SimpleAI{}}
	case cvc:
		playerX = ComputerPlayerWrapper{AI: SimpleAI{}}
		playerO = ComputerPlayerWrapper{AI: MinimaxAI{}}
	}
}

func drawMenu() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	width, height := termbox.Size()

	title := "TIC-TAC-TOE"
	drawText(width/2-len(title)/2, height/2-4, title, termbox.ColorWhite, termbox.ColorDefault)

	pvpText := "Player vs Player"
	pvcText := "Player vs Computer (Simple AI)"
	cvcText := "Computer vs Computer (Simple AI vs Minimax AI)"

	pvpColor := termbox.ColorDefault
	pvcColor := termbox.ColorDefault
	cvcColor := termbox.ColorDefault

	switch menuSelector {
	case 0:
		pvpColor = termbox.ColorGreen
	case 1:
		pvcColor = termbox.ColorGreen
	case 2:
		cvcColor = termbox.ColorGreen
	}

	drawText(width/2-len(pvpText)/2, height/2-1, pvpText, pvpColor, termbox.ColorDefault)
	drawText(width/2-len(pvcText)/2, height/2, pvcText, pvcColor, termbox.ColorDefault)
	drawText(width/2-len(cvcText)/2, height/2+1, cvcText, cvcColor, termbox.ColorDefault)

	instructions := "Use Arrow Keys to select, Enter to confirm"
	drawText(width/2-len(instructions)/2, height-2, instructions, termbox.ColorDefault, termbox.ColorDefault)

	termbox.Flush()
}

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	width, height := termbox.Size()

	// Instructions
	drawText(0, 0, "Use arrow keys to move, Enter to place, ESC for menu.", termbox.ColorDefault, termbox.ColorDefault)

	// Board template
	boardLines := []string{
		"   |   |   ",
		"---+---+---",
		"   |   |   ",
		"---+---+---",
		"   |   |   ",
	}

	startX := (width - len(boardLines[0])) / 2
	startY := (height - len(boardLines)) / 2

	// Draw the board template
	for i, line := range boardLines {
		drawText(startX, startY+i, line, termbox.ColorDefault, termbox.ColorDefault)
	}

	// Draw the X's and O's
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			if board[y][x] != ' ' {
				// Board coords (x,y) -> screen coords (startX+1 + x*4, startY + y*2)
				termbox.SetCell(startX+1+x*4, startY+y*2, board[y][x], termbox.ColorDefault, termbox.ColorDefault)
			}
		}
	}

	// Draw cursor
	if !gameOver && ((player == 'X' && playerX.IsHuman()) || (player == 'O' && playerO.IsHuman())) {
		// Cursor position needs to map to the same location as the X's and O's
		termbox.SetCursor(startX+1+cursorX*4, startY+cursorY*2)
	} else {
		termbox.HideCursor()
	}

	// Draw message
	messageY := startY + len(boardLines) + 1
	if message != "" {
		drawText(width/2-len(message)/2, messageY, message, termbox.ColorDefault, termbox.ColorDefault)
	} else {
		// Draw turn indicator
		turnMessage := ""
		if player == 'X' {
			turnMessage = fmt.Sprintf("%s's turn", playerX.Name())
		} else {
			turnMessage = fmt.Sprintf("%s's turn", playerO.Name())
		}
		drawText(width/2-len(turnMessage)/2, messageY, turnMessage, termbox.ColorDefault, termbox.ColorDefault)
	}

	termbox.Flush()
}

func drawText(x, y int, text string, fg, bg termbox.Attribute) {
	for i, r := range text {
		termbox.SetCell(x+i, y, r, fg, bg)
	}
}

func checkWin(currentBoard [boardHeight][boardWidth]rune, p rune) bool {
	// Check rows
	for y := 0; y < boardHeight; y++ {
		if currentBoard[y][0] == p && currentBoard[y][1] == p && currentBoard[y][2] == p {
			return true
		}
	}
	// Check columns
	for x := 0; x < boardWidth; x++ {
		if currentBoard[0][x] == p && currentBoard[1][x] == p && currentBoard[2][x] == p {
			return true
		}
	}
	// Check diagonals
	if currentBoard[0][0] == p && currentBoard[1][1] == p && currentBoard[2][2] == p {
		return true
	}
	if currentBoard[0][2] == p && currentBoard[1][1] == p && currentBoard[2][0] == p {
		return true
	}
	return false
}

func getEmptyCells(currentBoard [boardHeight][boardWidth]rune) [][2]int {
	var emptyCells [][2]int
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			if currentBoard[y][x] == ' ' {
				emptyCells = append(emptyCells, [2]int{x, y})
			}
		}
	}
	return emptyCells
}

func checkDraw(currentBoard [boardHeight][boardWidth]rune) bool {
	return len(getEmptyCells(currentBoard)) == 0 && !checkWin(currentBoard, 'X') && !checkWin(currentBoard, 'O')
}
