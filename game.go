package main

import (
    "fmt"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
    "image/color"
    "log"
    "sort"
)

const (
    screenWidth  = 320
    screenHeight = 240
    tileSize     = 5
)

type GameState int

const (
    StateMenu GameState = iota
    StateSelectDifficulty
    StatePlaying
    StateGameOver
    StateScoreboard
)

type Game struct {
    snake          *Snake
    food           *Food
    score          int
    gameOver       bool
    ticks          int
    updateCounter  float64
    speed          float64
    state          GameState
    highScores     []int
    speedIncrement float64
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
            g.state = StateSelectDifficulty
        } else if inpututil.IsKeyJustPressed(ebiten.KeyV) {
            g.state = StateScoreboard
        }
    case StateSelectDifficulty:
        if inpututil.IsKeyJustPressed(ebiten.Key1) {
            log.Println("Easy selected")
            g.startGame(10, 1)
        } else if inpututil.IsKeyJustPressed(ebiten.Key2) {
            log.Println("Medium selected")
            g.startGame(7, 3)
        } else if inpututil.IsKeyJustPressed(ebiten.Key3) {
            log.Println("Hard selected")
            g.startGame(5, 5)
        }
    case StatePlaying:
        g.updateGame()
    case StateGameOver:
        if inpututil.IsKeyJustPressed(ebiten.KeyR) {
            g.state = StateSelectDifficulty
        }
    case StateScoreboard:
        if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
            g.state = StateMenu
        }
    }
    return nil
}

func (g *Game) startGame(speed, increment float64) {
    g.snake = NewSnake()
    g.food = NewFood()
    g.score = 0
    g.speed = speed
    g.speedIncrement = increment
    g.gameOver = false
    g.state = StatePlaying
}

func (g *Game) updateGame() {
    g.updateCounter++
    if g.updateCounter < g.speed {
        return
    }
    g.updateCounter = 0

    g.snake.Move()

    if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.snake.Direction.X == 0 {
        g.snake.Direction = Point{X: -1, Y: 0}
    } else if ebiten.IsKeyPressed(ebiten.KeyRight) && g.snake.Direction.X == 0 {
        g.snake.Direction = Point{X: 1, Y: 0}
    } else if ebiten.IsKeyPressed(ebiten.KeyUp) && g.snake.Direction.Y == 0 {
        g.snake.Direction = Point{X: 0, Y: -1}
    } else if ebiten.IsKeyPressed(ebiten.KeyDown) && g.snake.Direction.Y == 0 {
        g.snake.Direction = Point{X: 0, Y: 1}
    }

    if g.collidesWithSnake() || g.collidesWithEdge() {
        g.gameOver = true
        g.state = StateGameOver
        g.highScores = append(g.highScores, g.score)
        sort.Slice(g.highScores, func(i, j int) bool {
            return g.highScores[i] > g.highScores[j]
        })
        if len(g.highScores) > 10 {
            g.highScores = g.highScores[:10]
        }
        g.saveHighScores()
        return
    }

    if g.snake.Body[0].X == g.food.Position.X && g.snake.Body[0].Y == g.food.Position.Y {
        g.snake.GrowCounter += 2
        g.score++
        g.speed -= g.speedIncrement
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
    case StateSelectDifficulty:
        g.DrawSelectDifficulty(screen)
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
    for _, segment := range g.snake.Body {
        x, y := segment.X*tileSize, segment.Y*tileSize
        ebitenutil.DrawRect(screen, float64(x), float64(y), tileSize, tileSize, color.White)
    }

    x, y := g.food.Position.X*tileSize, g.food.Position.Y*tileSize
    ebitenutil.DrawRect(screen, float64(x), float64(y), tileSize, tileSize, color.RGBA{R: 255, G: 0, B: 0, A: 255})

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
