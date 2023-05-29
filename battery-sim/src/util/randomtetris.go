package util

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

const (
	width       = 10 // Width of the game board
	height      = 20 // Height of the game board
	blockWidth  = 2  // Width of a Tetris block
	blockHeight = 1  // Height of a Tetris block
	startX      = 4  // Starting X coordinate for a new block
	startY      = 0  // Starting Y coordinate for a new block
)

type point struct {
	x, y int
}

type block struct {
	shape  [][]bool
	color  string
	offset point
}

type game struct {
	board     [][]string
	currBlock block
	nextBlock block
	score     int
	gameOver  bool
}

var blockShapes = [][][]bool{
	{
		{true, true},
		{true, true},
	},
	{
		{true, true, true, true},
	},
	{
		{true, true, true},
		{false, false, true},
	},
	{
		{true, true, true},
		{false, true, false},
	},
	{
		{true, true, true},
		{true, false, false},
	},
}

var blockColors = []string{
	"\033[31m", // Red
	"\033[32m", // Green
	"\033[33m", // Yellow
	"\033[34m", // Blue
	"\033[35m", // Magenta
}

func (g *game) init() {
	g.board = make([][]string, height)
	for i := range g.board {
		g.board[i] = make([]string, width)
	}

	g.currBlock = g.generateRandomBlock()
	g.nextBlock = g.generateRandomBlock()
	g.score = 0
	g.gameOver = false
}

func (g *game) generateRandomBlock() block {
	rand.Seed(time.Now().UnixNano())
	shape := blockShapes[rand.Intn(len(blockShapes))]
	color := blockColors[rand.Intn(len(blockColors))]
	offset := point{startX, startY}
	return block{shape: shape, color: color, offset: offset}
}

func (g *game) draw() {
	cmd := exec.Command("clear") // Use "cls" instead of "clear" on Windows
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Draw the board
	for _, row := range g.board {
		for _, cell := range row {
			if cell == "" {
				fmt.Print(" ")
			} else {
				fmt.Print(cell)
			}
		}
		fmt.Println()
	}

	fmt.Println("Score:", g.score)
	fmt.Println("Next Block:")

	// Draw the next block preview
	for _, row := range g.nextBlock.shape {
		for _, cell := range row {
			if cell {
				fmt.Print(g.nextBlock.color + "■")
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Println()
	}
}

func (g *game) canMove(x, y int) bool {
	for i := 0; i < blockHeight; i++ {
		for j := 0; j < blockWidth; j++ {
			if g.currBlock.shape[i][j] {
				newX := g.currBlock.offset.x + j + x
				newY := g.currBlock.offset.y + i + y

				// Check if the block is within the board boundaries
				if newX < 0 || newX >= width || newY >= height {
					return false
				}

				// Check if the block overlaps with existing blocks on the board
				if newY >= 0 && g.board[newY][newX] != "" {
					return false
				}
			}
		}
	}
	return true
}

func (g *game) move(x, y int) {
	if g.canMove(x, y) {
		g.currBlock.offset.x += x
		g.currBlock.offset.y += y
	}
}

func (g *game) rotate() {
	rotatedShape := make([][]bool, blockWidth)
	for i := range rotatedShape {
		rotatedShape[i] = make([]bool, blockHeight)
	}

	for i := 0; i < blockHeight; i++ {
		for j := 0; j < blockWidth; j++ {
			rotatedShape[j][blockHeight-i-1] = g.currBlock.shape[i][j]
		}
	}

	if g.canMove(0, 0) {
		g.currBlock.shape = rotatedShape
	}
}

func (g *game) placeBlock() {
	for i := 0; i < blockHeight; i++ {
		for j := 0; j < blockWidth; j++ {
			if g.currBlock.shape[i][j] {
				x := g.currBlock.offset.x + j
				y := g.currBlock.offset.y + i
				g.board[y][x] = g.currBlock.color + "■"
			}
		}
	}
}

func (g *game) clearLines() {
	fullLines := make([]int, 0)

	for i := 0; i < height; i++ {
		isFull := true
		for j := 0; j < width; j++ {
			if g.board[i][j] == "" {
				isFull = false
				break
			}
		}
		if isFull {
			fullLines = append(fullLines, i)
		}
	}

	for _, line := range fullLines {
		copy(g.board[1:line+1], g.board[:line])
		g.board[0] = make([]string, width)
		g.score++
	}
}

func (g *game) checkGameOver() {
	for i := 0; i < blockHeight; i++ {
		for j := 0; j < blockWidth; j++ {
			if g.currBlock.shape[i][j] {
				x := g.currBlock.offset.x + j
				y := g.currBlock.offset.y + i
				if y >= 0 && g.board[y][x] != "" {
					g.gameOver = true
					return
				}
			}
		}
	}
}

func main() {
	g := game{}
	g.init()

	for !g.gameOver {
		g.draw()
		g.checkGameOver()

		if !g.gameOver {
			g.placeBlock()
			g.clearLines()
			g.currBlock = g.nextBlock
			g.nextBlock = g.generateRandomBlock()

			g.getInput()
			g.update()
		}

		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("Game Over! Your score:", g.score)
}

func (g *game) getInput() {
	var input string
	fmt.Print("Enter move (left/right/rotate): ")
	fmt.Scan(&input)

	switch input {
	case "left":
		g.move(-1, 0)
	case "right":
		g.move(1, 0)
	case "rotate":
		g.rotate()
	default:
		fmt.Println("Invalid move. Please try again.")
		g.getInput()
	}
}

func (g *game) update() {
	g.move(0, 1)
}
