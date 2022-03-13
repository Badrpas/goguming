package foight

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"time"
)

type Entity struct {
	id   uint32
	game *Game

	x, y  float64
	angle float64

	body  *cp.Body
	shape *cp.Shape

	img          *ebiten.Image
	draw_options *ebiten.DrawImageOptions

	created_at, lifespan int64

	preupdate func(e *Entity, dt float64)
	update    func(e *Entity, dt float64)
	render    func(e *Entity, screen *ebiten.Image)

	on_dmg_received func(from *Entity, amount int32)
}

var id_counter uint32 = 0

func NewEntity(
	game *Game,
	x, y float64,
	body *cp.Body,
	shape *cp.Shape,
	image *ebiten.Image,
) *Entity {
	id_counter += 1

	return &Entity{
		id_counter,
		game,
		x, y, 0,
		body,
		shape,
		image,
		&ebiten.DrawImageOptions{},
		time.Now().UnixMilli(),
		-1,
		nil,
		nil,
		nil,
		nil,
	}

}

func (e *Entity) Update(dt float64) {
	if e.lifespan > 0 && (TimeNow()-e.created_at) >= e.lifespan {
		e.game.RemoveEntity(e)
		return
	}

	position := e.body.Position()
	e.x = position.X
	e.y = position.Y

	e.draw_options.GeoM.Reset()

	img_bounds := e.img.Bounds()
	e.draw_options.GeoM.Translate(
		float64(img_bounds.Dx()/-2),
		float64(img_bounds.Dy()/-2),
	)

	e.draw_options.GeoM.Rotate(e.angle)

	e.draw_options.GeoM.Translate(float64(e.x), float64(e.y))
}

func (e *Entity) Render(screen *ebiten.Image) {
	screen.DrawImage(e.img, e.draw_options)
}

func (e *Entity) RemoveFromGame() {
	e.game.RemoveEntity(e)
}
