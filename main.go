package main

import (
	"flag"
	levelmap "game/foight/map"
	"github.com/hajimehoshi/ebiten/v2"

	"log"

	"game/foight"
)

var mapname = flag.String("level", "levels/entry.tmx", "Map to run with")

func main() {

	g := foight.NewGame()

	err := levelmap.LoadToGameTiled(*mapname, g)
	if err != nil {
		log.Println("Couldn't load level", mapname)
		return
	}

	go foight.RunApi(g)

	ebiten.SetWindowSize(foight.ScreenWidth, foight.ScreenHeight)
	ebiten.SetWindowTitle("goguming")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
