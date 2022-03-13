package foight

import (
	imagestore "game/img"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/jakecoffman/cp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	clr "image/color"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"
)

var (
	mplusNormalFont font.Face
)

const (
	DEFAULT_HP         = 5
	INVINCIBILITY_TIME = int64(1000)
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Player struct {
	Entity

	game *Game

	name  string
	color clr.Color

	hp            int32
	is_invincible bool

	dx, dy float64
	tx, ty float64

	speed float64

	cooldown       int64
	last_fire_time int64

	messages chan UpdateMessage
}

func (g *Game) AddPlayer(name string, color clr.Color) *Player {

	player := &Player{
		game: g,

		name: name,
		hp:   DEFAULT_HP,

		Entity: *NewEntity(
			g,
			100+rand.Float64()*(ScreenWidth-200),
			100+rand.Float64()*(ScreenHeight-200),
			nil,
			nil,

			imagestore.Images["ploier.png"],
		),
		speed: 1000,

		cooldown:       300,
		last_fire_time: time.Now().UnixMilli(),

		messages: make(chan UpdateMessage, 1024),
	}

	player.preupdate = func(e *Entity, dt float64) {
		player.UpdateInputs(dt)
	}
	player.update = func(e *Entity, dt float64) {
		player.Update(dt)
	}
	player.render = func(e *Entity, screen *ebiten.Image) {
		e.Render(screen)

		info_hp := strings.Repeat("■", int(player.hp))
		l := float64(len(player.name))
		ll := float64(len(info_hp))

		text.Draw(screen, player.name, mplusNormalFont, int(player.x-l*4), int(player.y-62), player.color)
		text.Draw(screen, info_hp, mplusNormalFont, int(player.x-ll*2.5), int(player.y-42), player.color)
	}

	player.on_dmg_received = func(from *Entity, dmg int32) {
		if player.is_invincible {
			return
		}

		player.hp -= dmg

		if player.hp <= 0 {
			player.hp = DEFAULT_HP

			player.x = 100 + rand.Float64()*(ScreenWidth-200)
			player.y = 100 + rand.Float64()*(ScreenHeight-200)
			player.body.SetPosition(cp.Vector{player.x, player.y})
		}
	}

	player.SetColor(color)
	player.SetInvincible(INVINCIBILITY_TIME)

	g.AddEntity(&player.Entity)

	{ // Physics
		body := g.space.AddBody(cp.NewBody(1, 1))
		body.SetPosition(cp.Vector{player.x, player.y})
		body.UserData = &player.Entity

		radius := float64(imagestore.Images["ploier.png"].Bounds().Dx() / 2)
		shape := g.space.AddShape(cp.NewCircle(body, radius, cp.Vector{}))
		shape.SetElasticity(0.3)
		shape.SetFriction(0)
		shape.SetCollisionType(1)

		idx := player.Entity.id
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

	impulse := cp.Vector{tx, ty}
	p.body.ApplyImpulseAtLocalPoint(impulse, cp.Vector{})

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

func (p *Player) SetColor(color clr.Color) {
	p.color = color

	p.draw_options.ColorM.Scale(0, 0, 0, 1)

	rb, gb, bb, _ := color.RGBA()
	r := float64(rb) / 0xFFFF
	g := float64(gb) / 0xFFFF
	b := float64(bb) / 0xFFFF

	p.draw_options.ColorM.Translate(r, g, b, 0)
}

func (p *Player) SetInvincible(duration int64) {
	p.is_invincible = true
	precolor := p.color

	p.draw_options.ColorM.Scale(0.4, 0.4, 0.4, 1)

	p.timeholder.SetTimeout(func() {
		p.is_invincible = false
		p.SetColor(precolor)
	}, duration)
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
	b.draw_options.ColorM = p.draw_options.ColorM
	b.lifespan = 800
	b.on_dmg_dealt = on_bullet_dmg_dealt

	p.game.AddEntity(&b.Entity)

	dir := cp.Vector{p.tx, p.ty}.Normalize()

	p.body.Force()

	b.body.SetVelocityVector(dir.Mult(1000.).Add(p.body.Velocity().Mult(1.5)))
}

func on_bullet_dmg_dealt(b *Bullet, to *Entity) {
	b.Entity.RemoveFromGame()
}
