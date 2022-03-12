package foight

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"log"

	"image/color"
)

type Game struct {
	entities []*Entity

	space *cp.Space
}

func NewGame() *Game {
	game := &Game{
		space: cp.NewSpace(),
	}

	return game
}

func (g *Game) Layout(outWidth, outHeight int) (width, height int) {
	return 800, 600
}

func (g *Game) Update() error {
	var dt = 1. / 60. // Really disliking that

	for _, e := range g.entities {
		if e.preupdate != nil {
			e.preupdate(e, dt)
		}
	}

	g.space.Step(dt)

	for _, e := range g.entities {
		if e != nil {
			continue
		}

		if e.update != nil {
			e.update(e, dt)
		} else {
			e.Update(dt)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for _, e := range g.entities {
		e.Render(screen)
	}
}

func (g *Game) indexOfEntity(e *Entity) int32 {
	for k, v := range g.entities {
		if e == v {
			return int32(k)
		}
	}
	return -1 // not found.
}

func (g *Game) AddEntity(e *Entity) int32 {
	var idx int32 = -1

	g.entities = append(g.entities, e)

	idx = g.indexOfEntity(e)
	if idx < 0 {
		log.Fatal("Got negative idx for newly added player")
	}

	return idx
}

func (g *Game) RemoveEntity(e *Entity) {
	idx := g.indexOfEntity(e)

	if idx > -1 {
		g.entities = append(g.entities[:idx], g.entities[idx+1:]...)
	}
}
