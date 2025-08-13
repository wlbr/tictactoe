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
)

var (
	board        [boardHeight][boardWidth]rune
	cursorX      = 0
	cursorY      = 0
	player       = 'X'
	gameOver     = false
	winner       = ' '
	message      = ""
	gameMode     = -1 // -1: menu, 0: pvp, 1: pvc
	menuSelector = 0
)

func main() {
	Configure()
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
			handleGameEvent()
		}
	}
}

func handleMenuEvent() {
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyArrowUp, termbox.KeyArrowDown:
			menuSelector = (menuSelector + 1) % 2
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

func handleGameEvent() {
	if player == 'O' && gameMode == pvc && !gameOver {
		computerMove()
		draw()
		return
	}

	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		if gameOver {
			if ev.Key == termbox.KeyEnter {
				gameMode = -1 // Go back to menu
				resetGame()
			} else if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
				gameMode = -1 // Go back to menu
				resetGame()
			}
			return
		}

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
			if board[cursorY][cursorX] == ' ' {
				board[cursorY][cursorX] = player
				endTurn()
			}
		case termbox.KeyEsc, termbox.KeyCtrlC:
			gameMode = -1 // Go back to menu
			resetGame()
		}
	case termbox.EventError:
		panic(ev.Err)
	}
}

func endTurn() {
	if checkWin(player) {
		winner = player
		gameOver = true
		if gameMode == pvc && player == 'O' {
			message = "Computer wins! Press Enter to continue."
		} else {
			message = fmt.Sprintf("Player %c wins! Press Enter to continue.", player)
		}
	} else if checkDraw() {
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

func computerMove() {
	time.Sleep(500 * time.Millisecond) // A small delay to simulate thinking

	// 1. Check for a winning move
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			if board[y][x] == ' ' {
				board[y][x] = 'O'
				if checkWin('O') {
					endTurn()
					return
				}
				board[y][x] = ' ' // backtrack
			}
		}
	}

	// 2. Block player's winning move
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			if board[y][x] == ' ' {
				board[y][x] = 'X'
				if checkWin('X') {
					board[y][x] = 'O' // Place O to block
					endTurn()
					return
				}
				board[y][x] = ' ' // backtrack
			}
		}
	}

	// 3. Pick a random empty spot
	var emptyCells [][2]int
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			if board[y][x] == ' ' {
				emptyCells = append(emptyCells, [2]int{x, y})
			}
		}
	}

	if len(emptyCells) > 0 {
		move := emptyCells[rand.Intn(len(emptyCells))]
		board[move[1]][move[0]] = 'O'
		endTurn()
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
}

func drawMenu() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	width, height := termbox.Size()

	title := "TIC-TAC-TOE"
	drawText(width/2-len(title)/2, height/2-3, title, termbox.ColorWhite, termbox.ColorDefault)

	pvpText := "Player vs Player"
	pvcText := "Player vs Computer"

	pvpColor := termbox.ColorDefault
	pvcColor := termbox.ColorDefault
	if menuSelector == 0 {
		pvpColor = termbox.ColorGreen
	} else {
		pvcColor = termbox.ColorGreen
	}

	drawText(width/2-len(pvpText)/2, height/2, pvpText, pvpColor, termbox.ColorDefault)
	drawText(width/2-len(pvcText)/2, height/2+1, pvcText, pvcColor, termbox.ColorDefault)

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
	if !gameOver {
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
		if gameMode == pvc && player == 'O' {
			turnMessage = "Computer's turn..."
		} else {
			turnMessage = fmt.Sprintf("Player %c's turn", player)
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

func checkWin(p rune) bool {
	// Check rows
	for y := 0; y < boardHeight; y++ {
		if board[y][0] == p && board[y][1] == p && board[y][2] == p {
			return true
		}
	}
	// Check columns
	for x := 0; x < boardWidth; x++ {
		if board[0][x] == p && board[1][x] == p && board[2][x] == p {
			return true
		}
	}
	// Check diagonals
	if board[0][0] == p && board[1][1] == p && board[2][2] == p {
		return true
	}
	if board[0][2] == p && board[1][1] == p && board[2][0] == p {
		return true
	}
	return false
}

func checkDraw() bool {
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			if board[y][x] == ' ' {
				return false
			}
		}
	}
	return !checkWin('X') && !checkWin('O')
}
