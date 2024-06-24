package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	width        = 20
	height       = 10
	frameDelayMs = 100
)

type point struct {
	x, y int
}

type game struct {
	snake     []point
	food      point
	direction string
	gameOver  bool
	score     int
}

func (g *game) init() {
	g.snake = []point{{x: width / 2, y: height / 2}}
	g.food = point{x: rand.Intn(width), y: rand.Intn(height)}
	g.direction = "right"
	g.gameOver = false
	g.score = 0
}

// draw renders the game state to the console
func (g *game) draw() {
	fmt.Print("\033[H\033[2J") // Clear the screen

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if x == g.snake[0].x && y == g.snake[0].y {
				fmt.Print("O")
			} else if x == g.food.x && y == g.food.y {
				fmt.Print("*")
			} else {
				tail := false
				for _, p := range g.snake[1:] {
					if x == p.x && y == p.y {
						tail = true
						break
					}
				}
				if tail {
					fmt.Print("o")
				} else {
					fmt.Print(" ")
				}
			}
		}
		fmt.Println()
	}

	fmt.Println("Score:", g.score)
}

// update updates the game state
func (g *game) update() {
	head := g.snake[0]
	var newHead point

	// Move the snake's head based on the current direction
	switch g.direction {
	case "up":
		newHead = point{x: head.x, y: head.y - 1}
	case "down":
		newHead = point{x: head.x, y: head.y + 1}
	case "left":
		newHead = point{x: head.x - 1, y: head.y}
	case "right":
		newHead = point{x: head.x + 1, y: head.y}
	}

	// Prepend the new head to the snake
	g.snake = append([]point{newHead}, g.snake...)

	// Check for collisions with food
	if newHead.x == g.food.x && newHead.y == g.food.y {
		g.score++
		g.food = point{x: rand.Intn(width), y: rand.Intn(height)}
	} else {
		// Check if the snake has any body segments before removing one
		if len(g.snake) > 1 {
			g.snake = g.snake[:len(g.snake)-1]
		}
	}

	// Check for collisions with the walls or itself
	if newHead.x < 0 || newHead.x >= width || newHead.y < 0 || newHead.y >= height {
		g.gameOver = true
	} else {
		for _, p := range g.snake[1:] {
			if newHead.x == p.x && newHead.y == p.y {
				g.gameOver = true
				break
			}
		}
	}
}

func main() {
	// Initialize the random number generator
	rand.Seed(time.Now().UnixNano())

	// Initialize the game
	g := game{}
	g.init()

	// Open the keyboard
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	// Main game loop
	for !g.gameOver {
		// Draw the game state
		g.draw()

		// Update the game state
		g.update()

		// Get the next direction input from the user
		char, _, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		// Update the direction based on the input
		switch char {
		case 'z':
			g.direction = "up"
		case 's':
			g.direction = "down"
		case 'q':
			g.direction = "left"
		case 'd':
			g.direction = "right"
		case 'a':
			g.gameOver = true
		}

		// Delay between frames
		time.Sleep(frameDelayMs * time.Millisecond)
	}

	// Game over
	fmt.Println("Game Over! Your score:", g.score) 
}
