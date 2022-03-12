package foight

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jakecoffman/cp"
	"log"
)

var img_bullet *ebiten.Image
var img_w_bullet, img_h_bullet float64

func init() {
	var err error
	img_bullet, _, err = ebitenutil.NewImageFromFile("boolit.png")
	if err != nil {
		log.Fatal(err)
	}

	img_w_bullet = float64(img_bullet.Bounds().Dx())
	img_h_bullet = float64(img_bullet.Bounds().Dy())
}

type Bullet struct {
	Entity
}

func (g *Game) NewBullet(x, y float64) *Bullet {
	space := g.space

	body := space.AddBody(cp.NewBody(10, cp.INFINITY))
	body.SetPosition(cp.Vector{x, y})

	shape := space.AddShape(cp.NewCircle(body, img_w_bullet/2, cp.Vector{}))
	shape.SetElasticity(0)
	shape.SetFriction(0)
	shape.SetCollisionType(1)

	b := &Bullet{
		Entity: *NewEntity(
			g,
			x, y,
			0,
			body,
			shape,
			img_bullet,
			&ebiten.DrawImageOptions{},
		),
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

	hit := false
	e.body.EachArbiter(func(arbiter *cp.Arbiter) {
		hit = true
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

	log.Println("Do dmg!")

}
