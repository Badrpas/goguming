package foight

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"log"

	"image/color"
)

const (
	ScreenWidth  = 1600
	ScreenHeight = 1000
)

type Game struct {
	entities []*Entity

	space *cp.Space
}

func NewGame() *Game {
	game := &Game{
		space: cp.NewSpace(),
	}

	game.space.Iterations = 10
	game.space.SetDamping(0.1)

	addWalls(game.space)

	return game
}

func addWalls(space *cp.Space) {

	walls := []cp.Vector{
		{0, 0}, {0, ScreenHeight},
		{ScreenWidth, 0}, {ScreenWidth, ScreenHeight},
		{0, 0}, {ScreenWidth, 0},
		{0, ScreenHeight}, {ScreenWidth, ScreenHeight},
	}

	for i := 0; i < len(walls)-1; i += 2 {
		shape := space.AddShape(cp.NewSegment(space.StaticBody, walls[i], walls[i+1], 10))
		shape.SetElasticity(1)
		shape.SetFriction(1)
	}
}

func (g *Game) Layout(outWidth, outHeight int) (width, height int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	dt := 1. / 60. // Really disliking that

	for _, e := range g.entities {
		if e.preupdate != nil {
			e.preupdate(e, dt)
		}
	}

	g.space.Step(dt)

	for _, e := range g.entities {
		if e == nil {
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
		if e.render != nil {
			e.render(e, screen)
		} else {
			e.Render(screen)
		}
	}
}

func (g *Game) indexOfEntity(e *Entity) int32 {
	for k, v := range g.entities {
		if e == v || e.id == v.id { // Should not be dependent on the id
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
	if idx == -1 {
		return
	}

	if e.body != nil {
		g.space.RemoveBody(e.body)
	}
	if e.shape != nil {
		g.space.RemoveShape(e.shape)
	}

	g.entities = append(g.entities[:idx], g.entities[idx+1:]...)
}
