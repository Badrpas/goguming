package foight

import (
	imagestore "game/img"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
)

var _BULLET_IMG *ebiten.Image
var img_w_bullet, img_h_bullet float64

func init() {
	_BULLET_IMG = imagestore.Images["boolit.png"]

	img_w_bullet = float64(_BULLET_IMG.Bounds().Dx())
	img_h_bullet = float64(_BULLET_IMG.Bounds().Dy())
}

type Bullet struct {
	Entity

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
		Entity: *NewEntity(
			g,
			x, y,
			body,
			shape,
			_BULLET_IMG,
		),

		dmg: 1,
	}

	body.UserData = &b.Entity

	b.Entity.UpdateFn = func(e *Entity, dt float64) {
		b.Update(dt)
	}

	return b
}

func (b *Bullet) Update(dt float64) {
	e := b.Entity
	e.Update(dt)

	e.Body.EachArbiter(func(arbiter *cp.Arbiter) {
		b1, b2 := arbiter.Bodies()
		b.applyDamageTo(b1.UserData)
		b.applyDamageTo(b2.UserData)
	})
}

func (b *Bullet) applyDamageTo(i interface{}) {
	var entity *Entity

	if e, ok := i.(*Entity); ok {
		entity = e
	} else {
		return
	}

	if entity == &b.Entity {
		return
	}

	if entity.OnDmgReceived != nil {
		entity.OnDmgReceived(&b.Entity, b.dmg)

		if b.on_dmg_dealt != nil {
			b.on_dmg_dealt(b, entity)
		}
	}

}
