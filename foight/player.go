package foight

import (
	"game/foight/mixins"
	"game/foight/net"
	imagestore "game/img"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/jakecoffman/cp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	imagecolor "image/color"
	"log"
	"math/rand"
	"time"
)

var (
	mplusNormalFont font.Face
	_PLAYER_IMAGE   = imagestore.Images["ploier.png"]
)

const (
	DEFAULT_HP         = 5
	INVINCIBILITY_TIME = int64(2.5 * 1000)
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
	*Unit

	mixins.KDA

	messages chan net.UpdateMessage
}

func NewPlayer(g *Game, name string, color imagecolor.Color) *Player {

	player := &Player{

		Unit: NewUnit(name, 200, 200, _PLAYER_IMAGE),

		messages: make(chan net.UpdateMessage, 1024),
	}

	player.Game = g
	player.Holder = player

	super_preupdate := player.PreUpdateFn
	player.PreUpdateFn = func(e *Entity, dt float64) {
		player.UpdateInputs(dt)
		super_preupdate(e, dt)
	}

	super_render := player.RenderFn
	player.RenderFn = func(e *Entity, screen *ebiten.Image) {
		super_render(e, screen)
		player.renderKda(screen)
	}

	player.onDeathFn = func(self *Unit) {
		player.Respawn()
		player.DeathCount++
	}

	player.SetColor(color)

	player.Team = player.ID
	player.Unit.Init(g)

	player.Respawn()

	return player
}

func (p *Player) renderKda(screen *ebiten.Image) {
	kda := p.KDA.ToString()
	lk := float64(len(kda))
	f := mplusNormalFont
	text.Draw(screen, kda, f, int(p.X-lk*4), int(p.Y+46), p.color)
}

func (player *Player) Respawn() {
	player.HP = DEFAULT_HP

	l := len(player.Game.PlayerSpawnPoints)
	if l > 0 {
		point := player.Game.PlayerSpawnPoints[rand.Int()%l]
		player.Body.SetPosition(point)
		player.X, player.Y = point.X, point.Y
	} else {
		player.X = 100 + rand.Float64()*(ScreenWidth-200)
		player.Y = 100 + rand.Float64()*(ScreenHeight-200)
		player.Body.SetPosition(cp.Vector{player.X, player.Y})
	}

	player.Body.SetVelocityVector(cp.Vector{})

	player.SetInvincible(INVINCIBILITY_TIME)
}

func (p *Player) UpdateInputs(dt float64) {
	p.readMessages()

	if p.is_fire_expected() {
		p.fire()
	}
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

func (p *Player) applyUpdateMessage(um *net.UpdateMessage) {
	p.Dx = float64(um.Dx) / 50
	p.Dy = float64(um.Dy) / 50
	p.Tx = float64(um.Tx) / 50
	p.Ty = float64(um.Ty) / 50
}

func (p *Player) is_fire_expected() bool {
	cooldownExpired := (time.Now().UnixMilli() - p.last_fire_time) > p.CoolDown
	triggerDown := (p.Tx*p.Tx + p.Ty*p.Ty) > 0.4
	return cooldownExpired && triggerDown
}

func (p *Player) fire() {
	p.last_fire_time = time.Now().UnixMilli()

	b := NewBullet(p.Game, p.X, p.Y)
	b.Shape.Filter.Group = p.Shape.Filter.Group
	b.DrawOpts.ColorM = p.DrawOpts.ColorM
	b.Lifespan = 800
	b.Issuer = p.Entity
	b.on_dmg_dealt = func(b *Bullet, to *Entity) {
		p.AttacksConnectedCount++
		b.Entity.RemoveFromGame()
	}

	p.Game.AddEntity(b.Entity)

	dir := cp.Vector{p.Tx, p.Ty}.Normalize()

	b.Body.SetVelocityVector(
		dir.Mult(p.ForceModifier * 600.).
			Add(p.Body.Velocity().Mult(1.5)),
	)
}
