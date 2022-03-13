package foight

import (
	imagestore "game/img"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
)

var img_bullet *ebiten.Image
var img_w_bullet, img_h_bullet float64

func init() {
	img_bullet = imagestore.Images["boolit.png"]

	img_w_bullet = float64(img_bullet.Bounds().Dx())
	img_h_bullet = float64(img_bullet.Bounds().Dy())
}

type Bullet struct {
	Entity

	dmg int32

	on_dmg_dealt func(b *Bullet, to *Entity)
}

func (g *Game) NewBullet(x, y float64) *Bullet {
	space := g.space

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
			img_bullet,
		),

		dmg: 1,
	}

	body.UserData = &b.Entity

	b.Entity.update = func(e *Entity, dt float64) {
		b.Update(dt)
	}

	return b
}

func (b *Bullet) Update(dt float64) {
	e := b.Entity
	e.Update(dt)

	e.body.EachArbiter(func(arbiter *cp.Arbiter) {
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

	if entity.on_dmg_received != nil {
		entity.on_dmg_received(&b.Entity, b.dmg)

		if b.on_dmg_dealt != nil {
			b.on_dmg_dealt(b, entity)
		}
	}

}
