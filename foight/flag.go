package foight

import (
	"fmt"
	imagestore "game/img"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/jakecoffman/cp"
)

type Flag struct {
	*Entity
}
type FlagHandler struct {
	*Entity

	Total, Collected int
}

func NewFlagHandler() *FlagHandler {
	handler := &FlagHandler{
		NewEntity(-10000, -10000, nil, nil, nil),
		0,
		0,
	}

	f := mplusNormalFont20

	handler.UpdateFn = func(e *Entity, dt float64) {
		if handler.Total == handler.Collected {
			handler.RenderFn = func(e *Entity, screen *ebiten.Image) {
				opts := &ebiten.DrawImageOptions{}
				x, y := e.Game.Camera.GetWorldCoords(20, 20)
				opts.GeoM.Scale(1.6, 1.6)
				opts.GeoM.Translate(x, y)
				e.Game.TranslateCamera(opts)
				str := "Collected all of 'em. Chisto krasivo"
				text.DrawWithOptions(screen, str, f, opts)
			}
		}
	}

	handler.RenderFn = func(e *Entity, screen *ebiten.Image) {
		opts := &ebiten.DrawImageOptions{}
		x, y := e.Game.Camera.GetWorldCoords(20, 60)
		opts.GeoM.Scale(1.6, 1.6)
		opts.GeoM.Translate(x, y)
		e.Game.TranslateCamera(opts)
		str := fmt.Sprintf("%d/%d Flags Collected", handler.Collected, handler.Total)
		text.DrawWithOptions(screen, str, f, opts)
	}

	return handler
}

func NewFlag(pos cp.Vector, handler *FlagHandler) *Flag {
	handler.Total++

	f := &Flag{
		NewEntity(pos.X, pos.Y, nil, nil, imagestore.Images["flag.png"]),
	}
	f.OnCollision = func(e, other *Entity) {
		_, ok := other.Holder.(*Player)
		if ok {
			handler.Collected++
			e.RemoveFromGame()
		}
	}
	return f
}

func (f *Flag) Init(game *Game) {
	body, shape := AddCirclePhysicsToEntity(game, f.Entity)

	body.SetType(cp.BODY_STATIC)

	shape.SetSensor(true)
	shape.SetCollisionType(1)
}
