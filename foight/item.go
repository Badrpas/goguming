package foight

import (
	imagestore "game/img"
	"github.com/jakecoffman/cp"
)

type Item struct {
	*Entity

	OnPickup func(player *Player)
}

func (i *Item) Init(game *Game) int32 {
	space := game.Space

	i.Body = space.AddBody(cp.NewBody(cp.INFINITY, cp.INFINITY))
	i.Body.SetPosition(cp.Vector{i.X, i.Y})
	i.Body.SetType(cp.BODY_STATIC)

	i.Shape = space.AddShape(cp.NewBox(i.Body, float64(i.Img.Bounds().Dx()), float64(i.Img.Bounds().Dy()), 0))
	i.Shape.SetElasticity(1)
	i.Shape.SetFriction(1)
	i.Shape.SetSensor(true)

	i.Body.UserData = i.Entity
	i.OnCollision = func(e, other *Entity) {
		if i.OnPickup == nil {
			return
		}
		player, ok := other.Holder.(*Player)
		if ok {
			i.OnPickup(player)
		}
	}

	return game.AddEntity(i.Entity)
}

func newItem(pos cp.Vector) *Item {
	return &Item{
		Entity: NewEntity(
			nil,
			pos.X,
			pos.Y,
			nil,
			nil,
			nil,
		),
	}
}

func NewItemHeal(pos cp.Vector) *Item {
	item := newItem(pos)
	item.Img = imagestore.Images["heal.png"]
	item.OnPickup = func(player *Player) {
		player.HP += 1
	}
	item.DrawOpts.ColorM.Scale(0, 0, 0, 1)
	item.DrawOpts.ColorM.Translate(0, 1, 0, 0)

	return item
}
