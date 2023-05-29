package ignorable

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

const (
	width  = 20 // Width of the game board
	height = 10 // Height of the game board
)

type point struct {
	x, y int
}

type game struct {
	snake     []point
	food      point
	direction string
	score     int
	gameOver  bool
}

func (g *game) init() {
	g.snake = []point{{width / 2, height / 2}}
	g.food = generateFood()
	g.direction = "right"
	g.score = 0
	g.gameOver = false
}

func generateFood() point {
	rand.Seed(time.Now().UnixNano())
	return point{rand.Intn(width), rand.Intn(height)}
}

func (g *game) draw() {
	cmd := exec.Command("clear") // Use "cls" instead of "clear" on Windows
	cmd.Stdout = os.Stdout
	cmd.Run()

	board := make([][]string, height)
	for i := range board {
		board[i] = make([]string, width)
	}

	// Draw snake
	for _, p := range g.snake {
		board[p.y][p.x] = "■"
	}

	// Draw food
	board[g.food.y][g.food.x] = "★"

	// Print board
	for _, row := range board {
		for _, cell := range row {
			fmt.Print(cell + " ")
		}
		fmt.Println()
	}

	fmt.Println("Score:", g.score)
}

func (g *game) getInput() {
	var input string
	fmt.Print("Enter direction (up/down/left/right): ")
	fmt.Scan(&input)

	switch input {
	case "up":
		g.direction = "up"
	case "down":
		g.direction = "down"
	case "left":
		g.direction = "left"
	case "right":
		g.direction = "right"
	default:
		fmt.Println("Invalid direction. Please try again.")
		g.getInput()
	}
}

func (g *game) update() {
	head := g.snake[0]
	var newHead point

	switch g.direction {
	case "up":
		newHead = point{head.x, head.y - 1}
	case "down":
		newHead = point{head.x, head.y + 1}
	case "left":
		newHead = point{head.x - 1, head.y}
	case "right":
		newHead = point{head.x + 1, head.y}
	}

	// Check if the snake hits the wall
	if newHead.x < 0 || newHead.x >= width || newHead.y < 0 || newHead.y >= height {
		g.gameOver = true
		return
	}

	// Check if the snake hits itself
	for _, p := range g.snake[1:] {
		if newHead == p {
			g.gameOver = true
			return
		}
	}

	// Check if the snake eats the food
	if newHead == g.food {
		g.score++
		g.food = generateFood()
	} else {
		g.snake = g.snake[:len(g.snake)-1]
	}

	g.snake = append([]point{newHead}, g.snake...)
}

func main() {
	g := game{}
	g.init()

	for !g.gameOver {
		g.draw()
		g.getInput()
		g.update()
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("Game Over! Your score:", g.score)
}
