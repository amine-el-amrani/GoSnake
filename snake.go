package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 320
	screenHeight = 240
	popUpWidth    = 500
	popUpHeight   = 500
	tileSize      = 5
	borderSize    = 10
	borderPadding = 5
)

type Point struct {
	X int
	Y int
}

type Snake struct {
	Body        []Point
	Direction   Point
	GrowCounter int
}

type Food struct {
	Position Point
}

type Game struct {
	snake         *Snake
	food          *Food
	score         int
	gameOver      bool
	ticks         int
	updateCounter float64
	speed         float64
}

func main() {
	rand.Seed(time.Now().UnixNano())
	game := &Game{
		snake:    NewSnake(),
		food:     NewFood(),
		gameOver: false,
		ticks:    0,
		speed:    10,
	}
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Snake Game")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func NewSnake() *Snake {
	return &Snake{
		Body: []Point{
			{X: screenWidth / tileSize / 2, Y: screenHeight / tileSize / 2},
		},
		Direction: Point{X: 1, Y: 0},
	}
}

func (s *Snake) Move() {
	newHead := Point{
		X: s.Body[0].X + s.Direction.X,
		Y: s.Body[0].Y + s.Direction.Y,
	}
	s.Body = append([]Point{newHead}, s.Body...)

	if s.GrowCounter > 0 {
		s.GrowCounter--
	} else {
		s.Body = s.Body[:len(s.Body)-1]
	}
}

func NewFood() *Food {
	return &Food{
		Position: Point{
			X: rand.Intn(screenWidth / tileSize),
			Y: rand.Intn(screenHeight / tileSize),
		},
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	if g.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.restart()
		}
		return nil
	}

	g.updateCounter++
	if g.updateCounter < g.speed {
		return nil
	}
	g.updateCounter = 0

	// Update the snake's position
	g.snake.Move()

	// Handle user input
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.snake.Direction.X == 0 {
		g.snake.Direction = Point{X: -1, Y: 0}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) && g.snake.Direction.X == 0 {
		g.snake.Direction = Point{X: 1, Y: 0}
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) && g.snake.Direction.Y == 0 {
		g.snake.Direction = Point{X: 0, Y: -1}
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) && g.snake.Direction.Y == 0 {
		g.snake.Direction = Point{X: 0, Y: 1}
	}

	// Check for collisions
	if g.collidesWithSnake() || g.collidesWithEdge() {
		g.gameOver = true
		return nil
	}

	// Update food position if snake eats it
	if g.snake.Body[0].X == g.food.Position.X && g.snake.Body[0].Y == g.food.Position.Y {
		g.snake.GrowCounter += 2
		g.score++
        g.speed -= 0.5
		g.food = NewFood()
	}

	return nil
}

func (g *Game) collidesWithSnake() bool {
	head := g.snake.Body[0]
	for _, segment := range g.snake.Body[1:] {
		if head.X == segment.X && head.Y == segment.Y {
			return true
		}
	}
	return false
}

func (g *Game) collidesWithEdge() bool {
	head := g.snake.Body[0]
	return head.X < 0 || head.X >= screenWidth/tileSize || head.Y < 0 || head.Y >= screenHeight/tileSize
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen
	screen.Fill(color.Black)

	// Draw snake
	for _, segment := range g.snake.Body {
		x, y := segment.X*tileSize, segment.Y*tileSize
		ebitenutil.DrawRect(screen, float64(x), float64(y), tileSize, tileSize, color.White)
	}

	// Draw food
	x, y := g.food.Position.X*tileSize, g.food.Position.Y*tileSize
	ebitenutil.DrawRect(screen, float64(x), float64(y), tileSize, tileSize, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	// Draw score
	text := fmt.Sprintf("Score: %d", g.score)
	textX, textY := 10, 20
	ebitenutil.DebugPrint(screen, text)

	// Draw game over text if the game is over
	if g.gameOver {
		gameOverText := "Game Over. Press 'R' to restart."
		textWidth := len(gameOverText) * 20
		textX = (screenWidth - textWidth) / 2 + 100
		textY = screenHeight / 2 
		ebitenutil.DebugPrintAt(screen, gameOverText, textX, textY)
	}
}

func (g *Game) restart() {
	g.snake = NewSnake()
	g.food = NewFood()
	g.score = 0
	g.speed = 10
	g.gameOver = false
}
