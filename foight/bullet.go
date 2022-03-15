package foight

import (
	imagestore "game/img"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"log"
)

var _BULLET_IMG *ebiten.Image
var img_w_bullet, img_h_bullet float64

func init() {
	_BULLET_IMG = imagestore.Images["boolit.png"]

	img_w_bullet = float64(_BULLET_IMG.Bounds().Dx())
	img_h_bullet = float64(_BULLET_IMG.Bounds().Dy())
}

type Bullet struct {
	*Entity

	dmg int32

	on_dmg_dealt func(b *Bullet, to *Entity)
}

func NewBullet(g *Game, x, y float64) *Bullet {
	space := g.Space

	body := space.AddBody(cp.NewBody(10, 10))
	body.SetPosition(cp.Vector{x, y})

	shape := space.AddShape(cp.NewCircle(body, img_w_bullet/2, cp.Vector{}))
	shape.SetElasticity(0.5)
	shape.SetFriction(0)
	shape.SetCollisionType(1)

	b := &Bullet{
		Entity: NewEntity(
			g,
			x, y,
			body,
			shape,
			_BULLET_IMG,
		),

		dmg: 1,
	}

	body.UserData = b.Entity

	b.Entity.UpdateFn = func(e *Entity, dt float64) {
		b.Update(dt)
	}

	b.OnCollision = func(e, other *Entity) {
		b, ok := e.Holder.(*Bullet)
		if !ok {
			log.Fatalln("Received non bullet entity")
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
