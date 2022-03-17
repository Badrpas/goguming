package foight

import (
	"game/foight/mixins"
	"game/foight/net"
	"game/foight/util"
	imagestore "game/img"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/jakecoffman/cp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	imagecolor "image/color"
	"log"
	"math"
	"math/rand"
	"strings"
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
	*Entity

	Name string

	HP            int32
	is_invincible bool

	Dx, Dy float64
	Tx, Ty float64

	Speed         float64
	ForceModifier float64

	CoolDown       int64
	last_fire_time int64

	Effects []*Effect

	stunned_until int64

	mixins.KDA

	messages chan net.UpdateMessage
}

func NewPlayer(g *Game, name string, color imagecolor.Color) *Player {

	player := &Player{

		Name: name,
		HP:   DEFAULT_HP,

		Entity: NewEntity(
			-200,
			-200,
			nil,
			nil,

			_PLAYER_IMAGE,
		),

		Speed:         1000,
		ForceModifier: 1,

		CoolDown:       300,
		last_fire_time: time.Now().UnixMilli(),

		messages: make(chan net.UpdateMessage, 1024),
	}

	player.Game = g

	player.PreUpdateFn = func(e *Entity, dt float64) {
		player.UpdateEffects(dt)
		player.UpdateInputs(dt)
		player.Entity.PreUpdate(dt)
	}
	player.UpdateFn = func(e *Entity, dt float64) {
		player.Update(dt)
	}
	player.RenderFn = func(e *Entity, screen *ebiten.Image) {
		e.Render(screen)

		info_hp := strings.Repeat("■", int(player.HP))
		l := float64(len(player.Name))
		ll := float64(len(info_hp))

		f := mplusNormalFont
		text.Draw(screen, player.Name, f, int(player.X-l*4), int(player.Y-62), player.color)
		text.Draw(screen, info_hp, f, int(player.X-ll*2.5), int(player.Y-42), player.color)

		kda := player.KDA.ToString()
		lk := float64(len(kda))
		text.Draw(screen, kda, f, int(player.X-lk*4), int(player.Y+46), player.color)
	}

	player.OnDmgReceived = func(from *Entity, dmg int32) {
		if player.is_invincible {
			return
		}

		player.HP -= dmg

		if player.HP <= 0 {
			player.Respawn()
			player.DeathCount++
			bullet, ok := from.Holder.(*Bullet)
			if !ok {
				return
			}

			player, ok := bullet.Issuer.Holder.(*Player)
			if !ok {
				return
			}

			player.KillCount++
		}
	}
	player.Entity.Holder = player

	player.SetColor(color)

	g.AddEntity(player.Entity)

	{ // Physics
		body := g.Space.AddBody(cp.NewBody(1, 1))
		body.UserData = player.Entity

		radius := float64(_PLAYER_IMAGE.Bounds().Dx() / 2)
		shape := g.Space.AddShape(cp.NewCircle(body, radius, cp.Vector{}))
		shape.SetElasticity(0.3)
		shape.SetFriction(0)
		shape.SetCollisionType(1)

		idx := player.Entity.ID
		shape.Filter.Group = uint(idx + 1)

		player.Body = body
		player.Shape = shape
	}

	player.Respawn()

	return player
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

	player.StunFor(INVINCIBILITY_TIME * 9 / 10)
	player.SetInvincible(INVINCIBILITY_TIME)
}

func (p *Player) UpdateInputs(dt float64) {
	p.readMessages()

	if p.IsStunned() {
		return
	}

	tx := p.Dx * dt * p.Speed
	ty := p.Dy * dt * p.Speed

	impulse := cp.Vector{tx, ty}
	p.Body.ApplyImpulseAtLocalPoint(impulse, cp.Vector{})

	if p.is_fire_expected() {
		p.fire()
	}
}

func (p *Player) Update(dt float64) {
	if p.Tx != 0 || p.Ty != 0 {
		p.Angle = math.Atan2(float64(p.Ty), float64(p.Tx)) + math.Pi/2
	} else if p.Dx != 0 || p.Dy != 0 {
		p.Angle = math.Atan2(float64(p.Dy), float64(p.Dx)) + math.Pi/2
	}

	p.Entity.Update(dt)
}

func (p *Player) SetInvincible(duration int64) {
	p.is_invincible = true
	precolor := p.color

	p.Shape.SetSensor(true)
	p.DrawOpts.ColorM.Scale(0.4, 0.4, 0.4, 1)

	iteration := 1
	interval_id := p.TimeManager.SetInterval(func() {
		iteration++
		if iteration%2 == 0 {
			p.DrawOpts.ColorM.Scale(0.4, 0.4, 0.4, 1)
		} else {
			p.SetColor(precolor)
		}
	}, 300)

	p.TimeManager.SetTimeout(func() {
		p.TimeManager.ClearInterval(interval_id)
		p.Shape.SetSensor(false)
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
		dir.Mult(p.ForceModifier * 1000.).
			Add(p.Body.Velocity().Mult(1.5)),
	)
}

func (player *Player) StunFor(ms int64) {
	player.stunned_until = util.TimeNow() + ms
}
func (player *Player) IsStunned() bool {
	return util.TimeNow() < player.stunned_until
}

func (player *Player) UpdateEffects(dt float64) {
	count := len(player.Effects)
	row_length := int(math.Round(math.Sqrt(float64(count))))

	for idx, effect := range player.Effects {
		effect.Update(dt)

		x := idx % row_length
		y := idx / row_length
		effect.X, effect.Y = float64(x-row_length/2)*16, float64(y-row_length/2)*16
	}
}

func (player *Player) AddEffect(effect *Effect) {
	player.Effects = append(player.Effects, effect)

	effect.Parent = player.Entity

	if effect.Game != player.Game {
		player.Game.AddEntity(effect.Entity)
	}
}

func (player *Player) RemoveEffect(effect *Effect) {
	if effect.OnCease != nil {
		effect.OnCease(effect)
	}

	for idx, x := range player.Effects {
		if x == effect {
			player.Effects = append(player.Effects[:idx], player.Effects[idx+1:]...)
			break
		}
	}

	if effect.Game != nil {
		effect.RemoveFromGame()
	}
	effect.Target = nil
	effect.Data = nil
}
