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

	body := space.AddBody(cp.NewBody(1, cp.INFINITY))
	body.SetPosition(cp.Vector{x, y})

	shape := space.AddShape(cp.NewCircle(body, img_w_bullet/2, cp.Vector{}))
	shape.SetElasticity(0)
	shape.SetFriction(0)
	shape.SetCollisionType(1)

	return &Bullet{
		Entity: Entity{
			x, y,
			0,
			body,
			shape,
			img_bullet,
			&ebiten.DrawImageOptions{},
			nil,
			nil,
		},
	}
}

//func (b *Bullet) Update(dt float64) {
//	b.Entity.Update(dt)
//
//}
