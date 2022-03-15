package foight

import (
	imagestore "game/img"
	"github.com/jakecoffman/cp"
	imagecolor "image/color"
	"log"
)

type ItemCtor func(pos cp.Vector) *Item

var ItemConstructors []ItemCtor

func init() {
	ItemConstructors = []ItemCtor{
		NewItemHeal,
		NewItemSpeed,
		NewItemCoolDown,
	}
}

type Item struct {
	*Entity

	OnPickup func(player *Player)
}

func newItem(pos cp.Vector) *Item {
	return &Item{
		Entity: NewEntity(
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
	i.Body.SetType(cp.BODY_STATIC)

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

func NewItemWithEffect(pos cp.Vector, imgname string, color imagecolor.Color, effect *Effect) *Item {
	if effect == nil {
		log.Fatalln("nil instead of Effect")
	}
	item := newItem(pos)
	item.Img = imagestore.Images[imgname]
	item.SetColor(color)

	item.OnPickup = func(player *Player) {
		effect.ApplyTo(player)
	}

	return item
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

var Effects map[string]*Effect

func init() {
	Effects = map[string]*Effect{

		"Speed": {
			Entity: &Entity{
				Img:      imagestore.Images["speedup.png"],
				Lifespan: 6000,
				Scale:    cp.Vector{0.5, 0.5},
				color:    imagecolor.RGBA{255, 255, 0, 1},
			},
			OnApply: func(e *Effect) {
				player := e.Target
				delta := player.Speed * 0.8
				player.Speed += delta

				e.Data = delta
			},
			OnCease: func(e *Effect) {
				e.Target.Speed -= e.Data.(float64)
			},
		},

		"CoolDown": {
			Entity: &Entity{
				Img:      imagestore.Images["cooldown.png"],
				Lifespan: 6000,
				Scale:    cp.Vector{0.5, 0.5},
				color:    imagecolor.RGBA{0, 55, 255, 1},
			},
			OnApply: func(e *Effect) {
				player := e.Target
				delta := player.CoolDown / 2
				player.CoolDown -= delta

				e.Data = delta
			},
			OnCease: func(e *Effect) {
				e.Target.CoolDown += e.Data.(int64)
			},
		},
	}
}

func NewItemSpeed(pos cp.Vector) *Item {
	return NewItemWithEffect(
		pos,
		"speedup.png",
		imagecolor.RGBA{255, 255, 0, 0},
		Effects["Speed"],
	)
}

func NewItemCoolDown(pos cp.Vector) *Item {
	return NewItemWithEffect(
		pos,
		"cooldown.png",
		imagecolor.RGBA{0, 55, 255, 1},
		Effects["CoolDown"],
	)
}
