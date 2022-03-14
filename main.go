package main

import (
	levelmap "game/foight/map"
	"github.com/hajimehoshi/ebiten/v2"

	"log"

	"game/foight"
)

func main() {

	g := foight.NewGame()
	//err := levelmap.LoadToGameLdtk("levels/hola.ldtk", g)
	err := levelmap.LoadToGameTiled("levels/lul.tmx", g)
	if err != nil {
		log.Println("Couldn't load level")
		return
	}

	go foight.RunApi(g)

	ebiten.SetWindowSize(foight.ScreenWidth, foight.ScreenHeight)
	ebiten.SetWindowTitle("goguming")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
