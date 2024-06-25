package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "log"
)

func main() {
    game := &Game{
        snake:    NewSnake(),
        food:     NewFood(),
        gameOver: false,
        ticks:    0,
        speed:    10,
        state:    StateMenu,
        highScores: []int{},
    }

    game.loadHighScores()

    ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
    ebiten.SetWindowTitle("Snake Game")
    if err := ebiten.RunGame(game); err != nil {
        if err.Error() != "quit" {
            log.Fatal(err)
        }
        game.saveHighScores()
        log.Println("Game exited with quit signal.")
    }
}
