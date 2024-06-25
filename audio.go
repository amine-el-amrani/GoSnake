package main

import (
    "bytes"
    "github.com/hajimehoshi/ebiten/v2/audio"
    "github.com/hajimehoshi/ebiten/v2/audio/mp3"
    "io/ioutil"
    "log"
)

var (
    audioContext  *audio.Context
    eatPlayer     *audio.Player
    gameOverPlayer *audio.Player
	backgroundPlayer *audio.Player
)

func initAudio() {
    var err error
    audioContext = audio.NewContext(44100)

    eatPlayer, err = loadMP3("audio/eat.mp3")
    if err != nil {
        log.Fatalf("Failed to load eat sound: %v", err)
    }
	eatPlayer.SetVolume(2)

    gameOverPlayer, err = loadMP3("audio/gameover.mp3")
    if err != nil {
        log.Fatalf("Failed to load game over sound: %v", err)
    }
	gameOverPlayer.SetVolume(2)

	backgroundPlayer, err = loadMP3("audio/background.mp3")
    if err != nil {
        log.Fatalf("Failed to load background music: %v", err)
    }

    backgroundPlayer.SetVolume(1)
    backgroundPlayer.Play()
}

func loadMP3(path string) (*audio.Player, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    stream, err := mp3.Decode(audioContext, bytes.NewReader(data))
    if err != nil {
        return nil, err
    }

    p, err := audio.NewPlayer(audioContext, stream)
    if err != nil {
        return nil, err
    }
    return p, nil
}

func playSound(player *audio.Player) {
    player.Rewind()
    player.Play()
}
