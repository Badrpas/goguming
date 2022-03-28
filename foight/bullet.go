package foight

import (
	imagestore "game/img"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

var _BULLET_IMG *ebiten.Image

func init() {
	_BULLET_IMG = imagestore.Images["boolit.png"]
}

type Bullet struct {
	*Entity

	dmg    int32
	Issuer *Entity

	on_dmg_dealt func(b *Bullet, to *Entity)
}

func NewBullet(g *Game, x, y float64) *Bullet {

	b := &Bullet{
		Entity: NewEntity(
			x, y,
			nil,
			nil,
			_BULLET_IMG,
		),

		dmg: 1,
	}

	b.Holder = b

	body, shape := AddCirclePhysicsToEntity(g, b.Entity)
	body.SetMass(10)
	body.SetMoment(10)

	shape.SetElasticity(1)
	shape.SetFriction(0)
	shape.SetCollisionType(1)

	b.Entity.UpdateFn = func(e *Entity, dt float64) {
		b.Update(dt)
	}

	b.OnCollision = func(e, other *Entity) {
		b, ok := e.Holder.(*Bullet)
		if !ok {
			log.Fatalln("Received non bullet entity")
			return
		}

		if other.Team == e.Team {
			return
		}

		b.applyDamageTo(other)
	}

	return b
}

func (b *Bullet) Update(dt float64) {
	e := b.Entity
	e.Update(dt)
}

func (b *Bullet) applyDamageTo(entity *Entity) {
	if entity == b.Entity {
		return
	}

	if entity.OnDmgReceived != nil {
		entity.OnDmgReceived(b.Entity, b.dmg)

		if b.on_dmg_dealt != nil {
			b.on_dmg_dealt(b, entity)
		}
	}

}
