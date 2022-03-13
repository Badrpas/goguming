package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"log"

	"game/foight"
)

func main() {

	g := foight.NewGame()

	go foight.RunApi(g)

	ebiten.SetWindowSize(foight.ScreenWidth, foight.ScreenHeight)
	ebiten.SetWindowTitle("goguming")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
