package main

import (
	"encoding/json"
    "fmt"
	"log"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
    "image/color"
    "math/rand"
    "sort"
)

const (
    screenWidth  = 320
    screenHeight = 240
    tileSize     = 5
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

type GameState int

const (
    StateMenu GameState = iota
    StatePlaying
    StateGameOver
    StateScoreboard
)

type Game struct {
    snake         *Snake
    food          *Food
    score         int
    gameOver      bool
    ticks         int
    updateCounter float64
    speed         float64
    state         GameState
    highScores    []int
}


func (g *Game) DrawMenu(screen *ebiten.Image) {
    menuText1 := "Press 'S' to Start"
    menuText2 := "Press 'V' to View High Scores"
    menuText3 := "Press 'E' to Exit"
    textWidth1 := len(menuText1) * 20 + 70
    textWidth2 := len(menuText2) * 20 - 80
    textWidth3 := len(menuText3) * 20 + 85
    x1 := (screenWidth*2 - textWidth1) / 2
    x2 := (screenWidth*2 - textWidth2) / 2
    x3 := (screenWidth*2 - textWidth3) / 2
    y1 := screenHeight*2/ 7
    y2 := screenHeight*2/ 5
    y3 := screenHeight*2/ 4
    ebitenutil.DebugPrintAt(screen, menuText1, x1, y1)
    ebitenutil.DebugPrintAt(screen, menuText2, x2, y2)
    ebitenutil.DebugPrintAt(screen, menuText3, x3, y3)
}

func (g *Game) DrawScoreboard(screen *ebiten.Image) {
    scoreboardText := "High Scores"
    textWidth := len(scoreboardText) * 20 + 70
    x := (screenWidth*2 - textWidth) / 2 - 55
    y := screenHeight*2 / 5
    ebitenutil.DebugPrintAt(screen, scoreboardText, x, y)

    for i, score := range g.highScores {
        scoreText := fmt.Sprintf("%d. %d", i+1, score)
        textWidth := len(scoreText) * 6
        x := (screenWidth*2 - textWidth) / 2
        y := screenHeight*2/4 + 20*(i+1)
        ebitenutil.DebugPrintAt(screen, scoreText, x, y)
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
    if inpututil.IsKeyJustPressed(ebiten.KeyE) {
        log.Println("E key pressed. Quitting the game.")
        return fmt.Errorf("quit")
    }

    switch g.state {
    case StateMenu:
        if inpututil.IsKeyJustPressed(ebiten.KeyS) {
            g.startGame()
        } else if inpututil.IsKeyJustPressed(ebiten.KeyV) {
            g.state = StateScoreboard
        }
    case StatePlaying:
        g.updateGame()
    case StateGameOver:
        if inpututil.IsKeyJustPressed(ebiten.KeyR) {
            g.restart()
        }
    case StateScoreboard:
        if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
            g.state = StateMenu
        }
    }
    return nil
}



func (g *Game) startGame() {
    g.snake = NewSnake()
    g.food = NewFood()
    g.score = 0
    g.speed = 10
    g.gameOver = false
    g.state = StatePlaying
}

func (g *Game) updateGame() {
    g.updateCounter++
    if g.updateCounter < g.speed {
        return
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
        g.state = StateGameOver
        g.highScores = append(g.highScores, g.score)
        // Sort scores in descending order
        sort.Slice(g.highScores, func(i, j int) bool {
            return g.highScores[i] > g.highScores[j]
        })
        // Keep only the top 10 scores
        if len(g.highScores) > 10 {
            g.highScores = g.highScores[:10]
        }
        return
    }

    // Update food position if snake eats it
    if g.snake.Body[0].X == g.food.Position.X && g.snake.Body[0].Y == g.food.Position.Y {
        g.snake.GrowCounter += 2
        g.score++
        g.speed -= 0.5
        g.food = NewFood()
    }
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
    screen.Fill(color.Black)

    switch g.state {
    case StateMenu:
        g.DrawMenu(screen)
    case StatePlaying:
        g.drawGame(screen)	
	case StateGameOver:
        g.drawGame(screen)
        gameOverText1 := "Game Over. Press 'R' to restart."
        gameOverText2 := "Press 'E' to Exit."
        textWidth1 := len(gameOverText1) * 20
        textWidth2 := len(gameOverText2) * 20
        textX1 := (screenWidth*2 - textWidth1) / 2 + 65
        textX2 := (screenWidth*2 - textWidth2) / 2 - 30
        textY1 := screenHeight / 3
        textY2 := screenHeight / 2
        ebitenutil.DebugPrintAt(screen, gameOverText1, textX1, textY1)
        ebitenutil.DebugPrintAt(screen, gameOverText2, textX2, textY2)
    case StateScoreboard:
        g.DrawScoreboard(screen)
    }
}

func (g *Game) drawGame(screen *ebiten.Image) {
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
    ebitenutil.DebugPrint(screen, text)
}

func (g *Game) restart() {
    g.snake = NewSnake()
    g.food = NewFood()
    g.score = 0
    g.speed = 10
    g.gameOver = false
	g.state = StatePlaying
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
