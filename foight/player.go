package foight

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jakecoffman/cp"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"time"
)

var img *ebiten.Image
var img_w, img_h float64

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("ploier.png")
	if err != nil {
		log.Fatal(err)
	}

	img_w = float64(img.Bounds().Dx())
	img_h = float64(img.Bounds().Dy())
}

type Player struct {
	Entity

	game *Game

	name  string
	color uint32

	dx, dy float64
	tx, ty float64

	speed float64

	cooldown       int64
	last_fire_time int64

	messages chan UpdateMessage
}

func (g *Game) AddPlayer(name string, color uint32) *Player {

	player := &Player{
		game: g,

		name:  name,
		color: color,

		Entity: *NewEntity(
			g,
			100+rand.Float64()*300,
			100+rand.Float64()*300,
			0,
			nil,
			nil,

			img,
			&ebiten.DrawImageOptions{},
		),
		speed: 1000,

		cooldown:       300,
		last_fire_time: time.Now().UnixMilli(),

		messages: make(chan UpdateMessage, 1024),
	}

	player.Entity.preupdate = func(e *Entity, dt float64) {
		player.UpdateInputs(dt)
	}
	player.Entity.update = func(e *Entity, dt float64) {
		player.Update(dt)
	}

	player.SetColor(color)

	idx := g.AddEntity(&player.Entity)

	{ // Physics
		body := g.space.AddBody(cp.NewBody(1, 1))
		body.SetPosition(cp.Vector{player.x, player.y})
		body.UserData = &player.Entity

		shape := g.space.AddShape(cp.NewCircle(body, img_w/2, cp.Vector{}))
		shape.SetElasticity(0.3)
		shape.SetFriction(0)
		shape.SetCollisionType(1)

		shape.Filter.Group = uint(idx + 1)

		player.body = body
		player.shape = shape
	}

	return player
}

func (p *Player) UpdateInputs(dt float64) {
	p.readMessages()

	tx := p.dx * dt * p.speed
	ty := p.dy * dt * p.speed

	//p.body.SetVelocity(tx, ty)
	p.body.ApplyImpulseAtLocalPoint(cp.Vector{tx, ty}, cp.Vector{})

	if p.is_fire_expected() {
		p.fire()
	}
}

func (p *Player) Update(dt float64) {
	if p.tx != 0 || p.ty != 0 {
		p.angle = math.Atan2(float64(p.ty), float64(p.tx)) + math.Pi/2
	} else if p.dx != 0 || p.dy != 0 {
		p.angle = math.Atan2(float64(p.dy), float64(p.dx)) + math.Pi/2
	}

	p.Entity.Update(dt)
}

func (p *Player) SetColor(color uint32) {

	p.draw_options.ColorM.Scale(0, 0, 0, 1)

	r := float64((color&0xFF0000)>>4) / 0xff
	g := float64((color&0x00FF00)>>2) / 0xff
	b := float64((color&0x0000FF)>>0) / 0xff

	p.draw_options.ColorM.Translate(r, g, b, 0)
}

func (p *Player) readMessages() {
	for {
		select {
		case um := <-p.messages:
			p.applyUpdateMessage(&um)
		default:
			return
		}
	}
}

func (p *Player) applyUpdateMessage(um *UpdateMessage) {
	p.dx = float64(um.dx) / 50
	p.dy = float64(um.dy) / 50
	p.tx = float64(um.tx) / 50
	p.ty = float64(um.ty) / 50
}

func (p *Player) is_fire_expected() bool {
	cooldownExpired := (time.Now().UnixMilli() - p.last_fire_time) > p.cooldown
	triggerDown := (p.tx*p.tx + p.ty*p.ty) > 0.7
	return cooldownExpired && triggerDown
}

func (p *Player) fire() {
	p.last_fire_time = time.Now().UnixMilli()

	b := p.game.NewBullet(p.x, p.y)
	b.shape.Filter.Group = p.shape.Filter.Group
	b.lifespan = 1500

	p.game.AddEntity(&b.Entity)

	dir := cp.Vector{p.tx, p.ty}.Normalize()

	b.body.SetVelocityVector(dir.Mult(1000.))
}
