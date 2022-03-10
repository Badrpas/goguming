package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"log"

	"game/foight"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

func main() {

	g := foight.NewGame()

	go foight.RunApi(g)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("goguming")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
