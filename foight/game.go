package foight

import (
	"flag"
	"game/foight/debug"
	"game/foight/pathfind"
	"game/foight/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	camera "github.com/melonfunction/ebiten-camera"
	"log"
	"math/rand"
	"os"
	"time"

	"image/color"
)

var init_local = flag.Bool("local", false, "Add local player")

var local_player_added bool

const (
	ScreenWidth  = 1600
	ScreenHeight = 1000
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type Game struct {
	Camera   *camera.Camera
	Entities []*Entity

	PlayerSpawnPoints []cp.Vector
	ItemSpawnPoints   []cp.Vector

	Space *cp.Space
	Nav   *pathfind.Nav

	TimerManager *util.TimeHolder

	queued_jobs    chan func()
	removal_locked bool
}

func NewGame() *Game {
	game := &Game{
		Space:        cp.NewSpace(),
		Nav:          pathfind.NewNav(100, 100),
		TimerManager: &util.TimeHolder{},
		queued_jobs:  make(chan func(), 1024),
	}

	game.Space.Iterations = 10
	game.Space.SetDamping(0.8)

	//addWalls(game.Space)
	initItemSpawner(game)
	initCamera(game)

	if *init_local {
		addLocalPlayer(game)
	}

	test_unit := NewUnit("[NPC] Kekius", 80, 70, _PLAYER_IMAGE)
	test_unit.Team = test_unit.ID
	//test_unit.Speed = 100
	test_unit.Init(game)
	AddNpcController(test_unit)

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

func (g *Game) QueueJob(job func() interface{}) chan interface{} {
	c := make(chan interface{})

	g.queued_jobs <- func() {
		c <- job()
		close(c)
	}

	return c
}

func (g *Game) QueueJobVoid(job func()) {
	g.queued_jobs <- job
}

func (g *Game) runQueue() {
	for {
		select {
		case job := <-g.queued_jobs:
			job()
		default:
			return
		}
	}
}

func initItemSpawner(game *Game) {
	const ITEM_LIFESPAN = 8000
	const ITEM_SPAWN_INTERVAL = 3500
	last_spawns := make([]int, ITEM_LIFESPAN/ITEM_SPAWN_INTERVAL+4)
	l := len(last_spawns)
	storage_idx := 0

	get_next_point_idx := func() int {
		point_count := len(game.PlayerSpawnPoints)
		spawn_idx := 0
	Outer:
		for i := 0; i < l*2; //goland:noinspection GoUnreachableCode
		i++ {
			spawn_idx = rand.Int() % point_count
			for _, used_idx := range last_spawns {
				if used_idx == spawn_idx {
					continue Outer
				}
			}

			storage_idx = (storage_idx + 1) % l
			last_spawns[storage_idx] = spawn_idx
			return spawn_idx
		}

		return spawn_idx
	}

	game.TimerManager.SetInterval(func() {
		if l > 0 {
			point := game.PlayerSpawnPoints[get_next_point_idx()]
			ctor := ItemConstructors[rand.Int()%len(ItemConstructors)]
			item := ctor(point)
			item.Init(game)
			item.Lifespan = ITEM_LIFESPAN
		}
	}, ITEM_SPAWN_INTERVAL)

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

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	// Camera
	{
		cam_delta := dt * 100
		if ebiten.IsKeyPressed(ebiten.KeyH) {
			g.Camera.MovePosition(-cam_delta, 0)
		}
		if ebiten.IsKeyPressed(ebiten.KeyL) {
			g.Camera.MovePosition(cam_delta, 0)
		}
		if ebiten.IsKeyPressed(ebiten.KeyK) {
			g.Camera.MovePosition(0, -cam_delta)
		}
		if ebiten.IsKeyPressed(ebiten.KeyJ) {
			g.Camera.MovePosition(0, cam_delta)
		}

		UpdateCamera(g, dt)
	}

	g.TimerManager.Update()

	for _, e := range g.Entities {
		if e == nil {
			continue
		}

		if e.PreUpdateFn != nil {
			e.PreUpdateFn(e, dt)
		} else {
			e.PreUpdate(dt)
		}
	}

	g.removal_locked = true
	g.Space.Step(dt)
	g.removal_locked = false

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

	g.runQueue()
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		addLocalPlayer(g)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.Fill(color.Black)
	s := g.Camera.Surface
	s.Clear()

	for _, e := range g.Entities {
		if e.RenderFn != nil {
			e.RenderFn(e, s)
		} else {
			e.Render(s)
		}
	}

	for _, point := range debug.Points {
		opts := &ebiten.DrawImageOptions{}
		w, h := _BULLET_IMG.Size()
		opts.GeoM.Translate(float64(w/-2), float64(h/-2))
		opts.GeoM.Translate(point.X, point.Y)
		g.TranslateCamera(opts)
		s.DrawImage(_BULLET_IMG, opts)
	}

	Blit(g.Camera, screen)
}

func initCamera(game *Game) {
	game.Camera = camera.NewCamera(ScreenWidth, ScreenHeight, ScreenWidth/2, ScreenHeight/2, 0, 1)
	SetZoom(game.Camera, 0.9)
}

func (g *Game) TranslateCamera(opts *ebiten.DrawImageOptions) {
	c := g.Camera
	w, h := c.Width, c.Height
	opts.GeoM.Translate(float64(w)/c.Scale/2, float64(h)/c.Scale/2)
	opts.GeoM.Translate(-c.X, -c.Y)
}

func (g *Game) indexOfEntity(e *Entity) int32 {
	for k, v := range g.Entities {
		if e == v /*|| e.ID == v.ID*/ { // Should not be dependent on the id
			return int32(k)
		}
	}
	return -1 // not found.
}

func (g *Game) AddEntity(e *Entity) int32 {
	var idx = g.indexOfEntity(e)
	// Duplication guard
	if idx != -1 {
		return idx
	}

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
	if g.removal_locked {
		g.QueueJobVoid(func() {
			g.RemoveEntity(e)
		})
		return
	}

	idx := g.indexOfEntity(e)
	if idx == -1 {
		return
	}
	e.Holder = nil

	e.RemovePhysics()

	g.Entities = append(g.Entities[:idx], g.Entities[idx+1:]...)

	if e.OnRemove != nil {
		e.OnRemove(e)
	}

	e.Game = nil
}

func addLocalPlayer(g *Game) {
	if local_player_added {
		return
	}
	local_player_added = true

	player := NewPlayer(g, "Local Yoba", color.Gray{255})
	super_preupdate := player.PreUpdateFn
	player.PreUpdateFn = func(e *Entity, dt float64) {
		var dx, dy, tx, ty float64

		if ebiten.IsKeyPressed(ebiten.KeyA) {
			dx += -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			dx += 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			dy += -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			dy += 1
		}

		if dx != 0 || dy != 0 {
			v := cp.Vector{dx, dy}.Normalize()
			dx, dy = v.X, v.Y
		}

		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			tx += -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			tx += 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			ty += -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			ty += 1
		}

		if tx != 0 || ty != 0 {
			v := cp.Vector{tx, ty}.Normalize()
			tx, ty = v.X, v.Y
		}

		p := player
		p.Dx, p.Dy, p.Tx, p.Ty = dx, dy, tx, ty

		super_preupdate(e, dt)
	}
}
