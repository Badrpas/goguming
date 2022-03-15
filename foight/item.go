package foight

import (
	imagestore "game/img"
	"github.com/jakecoffman/cp"
)

type ItemCtor func(pos cp.Vector) *Item

var ItemConstructors []ItemCtor

func init() {
	ItemConstructors = []ItemCtor{
		NewItemHeal,
		NewItemSpeed,
	}
}

type Item struct {
	*Entity

	OnPickup func(player *Player)
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

func (i *Item) Init(game *Game) int32 {
	space := game.Space

	i.Body = space.AddBody(cp.NewBody(cp.INFINITY, cp.INFINITY))
	i.Body.SetPosition(cp.Vector{i.X, i.Y})
	//i.Body.SetType(cp.BODY_STATIC)

	bounds := i.Img.Bounds()
	box := cp.NewBox(i.Body, float64(bounds.Dx()), float64(bounds.Dy()), 0)
	i.Shape = space.AddShape(box)
	i.Shape.SetElasticity(1)
	i.Shape.SetFriction(1)
	i.Shape.SetSensor(true)
	i.Shape.SetCollisionType(1)

	i.Body.UserData = i.Entity
	i.OnCollision = func(e, other *Entity) {
		if i.OnPickup == nil {
			return
		}
		player, ok := other.Holder.(*Player)
		if ok {
			i.OnPickup(player)
			game.RemoveEntity(i.Entity)
		}
	}

	return game.AddEntity(i.Entity)
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

func NewItemSpeed(pos cp.Vector) *Item {
	item := newItem(pos)
	item.Img = imagestore.Images["speedup.png"]
	item.OnPickup = func(player *Player) {
		delta := player.Speed * 0.5
		player.Speed += delta
		player.TimeManager.SetTimeout(func() {
			player.Speed -= delta
		}, 3000)
	}
	item.DrawOpts.ColorM.Scale(0, 0, 0, 1)
	item.DrawOpts.ColorM.Translate(0, 1, 0, 0)

	return item
}
