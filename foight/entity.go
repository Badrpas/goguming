package foight

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"time"
)

type Entity struct {
	ID   uint32
	Game *Game

	X, Y  float64
	Angle float64

	Body  *cp.Body
	Shape *cp.Shape

	Img      *ebiten.Image
	DrawOpts *ebiten.DrawImageOptions

	CreatedAt, Lifespan int64

	TimeManager *TimeHolder

	Holder interface{}

	PreUpdateFn func(e *Entity, dt float64)
	UpdateFn    func(e *Entity, dt float64)
	RenderFn    func(e *Entity, screen *ebiten.Image)

	OnCollision   func(e, other *Entity)
	OnDmgReceived func(from *Entity, amount int32)
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
		&TimeHolder{},
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	}

}
func (e *Entity) SetRenderFn(fn func(e *Entity, screen *ebiten.Image)) {
	e.RenderFn = fn
}

func (e *Entity) Update(dt float64) {
	e.TimeManager.Update()

	if e.Lifespan > 0 && (TimeNow()-e.CreatedAt) >= e.Lifespan {
		e.Game.RemoveEntity(e)
		return
	}

	position := e.Body.Position()
	e.X = position.X
	e.Y = position.Y

	e.DrawOpts.GeoM.Reset()

	img_bounds := e.Img.Bounds()
	e.DrawOpts.GeoM.Translate(
		float64(img_bounds.Dx()/-2),
		float64(img_bounds.Dy()/-2),
	)

	e.DrawOpts.GeoM.Rotate(e.Angle)

	e.DrawOpts.GeoM.Translate(float64(e.X), float64(e.Y))

	if e.OnCollision != nil {
		e.Body.EachArbiter(func(arbiter *cp.Arbiter) {
			b1, b2 := arbiter.Bodies()
			if b1.UserData != nil && b1 != e.Body {
				e.OnCollision(e, b1.UserData.(*Entity))
			}
			if b2.UserData != nil && b2 != e.Body {
				e.OnCollision(e, b2.UserData.(*Entity))
			}
		})
	}
}

func (e *Entity) Render(screen *ebiten.Image) {
	screen.DrawImage(e.Img, e.DrawOpts)
}

func (e *Entity) RemoveFromGame() {
	e.Game.RemoveEntity(e)
}
