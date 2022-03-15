package foight

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"log"
	"math/rand"
	"time"

	"image/color"
)

const (
	ScreenWidth  = 1600
	ScreenHeight = 1000
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type Game struct {
	Entities []*Entity

	PlayerSpawnPoints []cp.Vector
	ItemSpawnPoints   []cp.Vector

	Space *cp.Space

	TimerManager *TimeHolder
}

func NewGame() *Game {
	game := &Game{
		Space:        cp.NewSpace(),
		TimerManager: &TimeHolder{},
	}

	game.Space.Iterations = 10
	game.Space.SetDamping(0.1)

	addWalls(game.Space)
	initItemSpawner(game)

	handler := game.Space.NewCollisionHandler(1, 1)
	handler.BeginFunc = func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		b1, b2 := arb.Bodies()
		var e1, e2 *Entity
		if b1.UserData != nil {
			e1, _ = b1.UserData.(*Entity)
		}
		if b2.UserData != nil {
			e2, _ = b2.UserData.(*Entity)
		}

		if e1 != nil && e2 != nil {
			if e1.OnCollision != nil {
				e1.OnCollision(e1, e2)
			}
			if e2.OnCollision != nil {
				e2.OnCollision(e2, e1)
			}
		}

		return true
	}

	return game
}

func initItemSpawner(game *Game) {
	game.TimerManager.SetInterval(func() {
		l := len(game.PlayerSpawnPoints)
		if l > 0 {
			point := game.PlayerSpawnPoints[rand.Int()%l]
			ctor := ItemConstructors[rand.Int()%len(ItemConstructors)]
			item := ctor(point)
			item.Init(game)
			item.Lifespan = 20000
		}
	}, 10000)
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

	g.TimerManager.Update()

	for _, e := range g.Entities {
		if e.PreUpdateFn != nil {
			e.PreUpdateFn(e, dt)
		}
	}

	g.Space.Step(dt)

	for _, e := range g.Entities {
		if e == nil {
			continue
		}

		if e.UpdateFn != nil {
			e.UpdateFn(e, dt)
		} else {
			e.Update(dt)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for _, e := range g.Entities {
		if e.RenderFn != nil {
			e.RenderFn(e, screen)
		} else {
			e.Render(screen)
		}
	}
}

func (g *Game) indexOfEntity(e *Entity) int32 {
	for k, v := range g.Entities {
		if e == v || e.ID == v.ID { // Should not be dependent on the id
			return int32(k)
		}
	}
	return -1 // not found.
}

func (g *Game) AddEntity(e *Entity) int32 {
	var idx int32 = -1
	if e.Game != g {
		e.Game = g
	}

	g.Entities = append(g.Entities, e)

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
	e.Holder = nil

	if e.Body != nil {
		g.Space.RemoveBody(e.Body)
	}
	if e.Shape != nil {
		g.Space.RemoveShape(e.Shape)
	}

	g.Entities = append(g.Entities[:idx], g.Entities[idx+1:]...)
}
