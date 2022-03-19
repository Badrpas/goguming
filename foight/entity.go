package foight

import (
	"game/foight/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	imagecolor "image/color"
	"log"
)

type Entity struct {
	ID   uint32
	Game *Game

	X, Y  float64
	Angle float64

	Parent *Entity

	Body  *cp.Body
	Shape *cp.Shape

	color imagecolor.Color

	Img      *ebiten.Image
	DrawOpts *ebiten.DrawImageOptions
	Scale    cp.Vector

	CreatedAt, Lifespan int64

	TimeManager *util.TimeHolder

	Holder interface{}

	PreUpdateFn func(e *Entity, dt float64)
	UpdateFn    func(e *Entity, dt float64)
	RenderFn    func(e *Entity, screen *ebiten.Image)

	OnCollision   func(e, other *Entity)
	OnDmgReceived func(from *Entity, amount int32)

	OnRemove func(entity *Entity)
}

var id_counter uint32 = 0

func NewEntity(
	x, y float64,
	body *cp.Body,
	shape *cp.Shape,
	image *ebiten.Image,
) *Entity {
	id_counter += 1

	return &Entity{
		id_counter,
		nil,
		x, y, 0,
		nil,
		body,
		shape,
		imagecolor.White,
		image,
		&ebiten.DrawImageOptions{},
		cp.Vector{1, 1},
		util.TimeNow(),
		-1,
		&util.TimeHolder{},
		nil,
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

func (e *Entity) PreUpdate(dt float64) {
	if e.Parent != nil && e.Body != nil {
		if e.Parent.Body != nil {
			e.Body.SetPosition(e.Parent.Body.Position().Add(cp.Vector{e.X, e.Y}))
		} else {
			e.Body.SetPosition(cp.Vector{e.Parent.X + e.X, e.Parent.Y + e.Y})
		}
	}
}

func (e *Entity) Update(dt float64) {
	if e.Game == nil {
		log.Println("Got Entity.Update() without .Game being present")
		return
	}
	e.TimeManager.Update()

	if e.Lifespan > 0 && (util.TimeNow()-e.CreatedAt) >= e.Lifespan {
		// Should it be handled differently?
		if e.Game != nil {
			e.Game.RemoveEntity(e)
		}
		return
	}

	e.DrawOpts.GeoM.Reset()

	img_bounds := e.Img.Bounds()
	e.DrawOpts.GeoM.Translate(
		float64(img_bounds.Dx()/-2),
		float64(img_bounds.Dy()/-2),
	)
	e.DrawOpts.GeoM.Scale(e.Scale.X, e.Scale.Y)

	e.DrawOpts.GeoM.Rotate(e.Angle)

	if e.Body == nil {
		if e.Parent != nil {
			e.DrawOpts.GeoM.Translate(float64(e.X)+e.Parent.X, float64(e.Y)+e.Parent.Y)
			return
		}
	} else {
		position := e.Body.Position()
		e.X = position.X
		e.Y = position.Y
	}

	//e.DrawOpts.GeoM.Translate(float64(e.X), float64(e.Y))
	e.translateCamera()
}

func (e *Entity) translateCamera() {
	c := e.Game.Camera
	w, h := c.Surface.Size()
	e.DrawOpts.GeoM.Translate(float64(w)/2, float64(h)/2)
	e.DrawOpts.GeoM.Translate(-c.X+e.X, -c.Y+e.Y)
}

func (e *Entity) Render(screen *ebiten.Image) {
	screen.DrawImage(e.Img, e.DrawOpts)
}

func (e *Entity) SetColor(color imagecolor.Color) {
	if color == nil {
		return
	}
	e.color = color

	e.DrawOpts.ColorM.Scale(0, 0, 0, 1)

	rb, gb, bb, _ := color.RGBA()
	r := float64(rb) / 0xFFFF
	g := float64(gb) / 0xFFFF
	b := float64(bb) / 0xFFFF

	e.DrawOpts.ColorM.Translate(r, g, b, 0)
}

func (e *Entity) RemoveFromGame() {
	e.Game.RemoveEntity(e)
}

func (e *Entity) RemovePhysics() {
	if e.Body != nil {
		e.Game.Space.RemoveBody(e.Body)
		e.Body.UserData = nil
		e.Body = nil
	}
	if e.Shape != nil {
		e.Game.Space.RemoveShape(e.Shape)
		e.Shape = nil
	}
}
