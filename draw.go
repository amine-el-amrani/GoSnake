package main

import (
	"fmt"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

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

func (g *Game) DrawSelectDifficulty(screen *ebiten.Image) {
    difficultyText1 := "Press '1' for Easy"
    difficultyText2 := "Press '2' for Medium"
    difficultyText3 := "Press '3' for Hard"
    textWidth1 := len(difficultyText1) * 20 + 70
    textWidth2 := len(difficultyText2) * 20 + 40
    textWidth3 := len(difficultyText3) * 20 + 70
    x1 := (screenWidth*2 - textWidth1) / 2
    x2 := (screenWidth*2 - textWidth2) / 2
    x3 := (screenWidth*2 - textWidth3) / 2
    y1 := screenHeight*2 / 7
    y2 := screenHeight*2 / 5
    y3 := screenHeight*2 / 4
    ebitenutil.DebugPrintAt(screen, difficultyText1, x1, y1)
    ebitenutil.DebugPrintAt(screen, difficultyText2, x2, y2)
    ebitenutil.DebugPrintAt(screen, difficultyText3, x3, y3)
}

func (g *Game) DrawScoreboard(screen *ebiten.Image) {
    scoreboardText := "High Scores"
    textWidth := len(scoreboardText) * 20 + 70
    x := (screenWidth*2 - textWidth) / 2 - 55
    y := screenHeight*2 / 10
    ebitenutil.DebugPrintAt(screen, scoreboardText, x, y)

    for i, score := range g.highScores {
        scoreText := fmt.Sprintf("%d. %d", i+1, score)
        textWidth := len(scoreText) * 20 + 70
        x := (screenWidth*2 - textWidth) / 2 - 105
        y := screenHeight*2/8 + 20*(i+1)
        ebitenutil.DebugPrintAt(screen, scoreText, x, y)
    }
}
