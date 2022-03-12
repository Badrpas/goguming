package foight

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
)

type Entity struct {
	x, y float64

	angle float64

	body  *cp.Body
	shape *cp.Shape

	img          *ebiten.Image
	draw_options *ebiten.DrawImageOptions

	preupdate func(e *Entity, dt float64)
	update    func(e *Entity, dt float64)
}

func (e *Entity) Update(dt float64) {
	position := e.body.Position()
	e.x = position.X
	e.y = position.Y

	e.draw_options.GeoM.Reset()
	e.draw_options.GeoM.Translate(img_w/-2, img_h/-2)

	e.draw_options.GeoM.Rotate(e.angle)

	e.draw_options.GeoM.Translate(float64(e.x), float64(e.y))
}

func (e *Entity) Render(screen *ebiten.Image) {
	screen.DrawImage(e.img, e.draw_options)
}