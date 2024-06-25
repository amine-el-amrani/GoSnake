package main

import (
    "encoding/json"
    "log"
    "os"
)

func (g *Game) saveHighScores() {
    file, err := os.Create("highscores.json")
    if err != nil {
        log.Println("Error saving high scores:", err)
        return
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    err = encoder.Encode(g.highScores)
    if err != nil {
        log.Println("Error encoding high scores:", err)
    }
}

func (g *Game) loadHighScores() {
    file, err := os.Open("highscores.json")
    if err != nil {
        if !os.IsNotExist(err) {
            log.Println("Error loading high scores:", err)
        }
        return
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    err = decoder.Decode(&g.highScores)
    if err != nil {
        log.Println("Error decoding high scores:", err)
    }
}
