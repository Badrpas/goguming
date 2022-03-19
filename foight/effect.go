package foight

import "github.com/hajimehoshi/ebiten/v2"

type EffectCallback func(e *Effect)

type Effect struct {
	*Entity

	Target *Player
	Data   interface{}

	OnApply EffectCallback
	OnCease EffectCallback
}

// Returns newly applied (cloned) Effect
func (proto *Effect) ApplyTo(player *Player) *Effect {
	_e := *proto
	effect := &_e

	effect.Entity = NewEntity(0, 0, nil, nil, effect.Img)
	effect.Scale = proto.Scale

	effect.Target = player
	effect.Lifespan = proto.Lifespan
	effect.SetColor(proto.color)

	effect.DrawOpts.Filter = ebiten.FilterLinear

	effect.UpdateFn = func(en *Entity, dt float64) {
		effect.Update(dt)
	}
	effect.OnRemove = func(entity *Entity) {
		player.RemoveEffect(effect)
	}

	player.AddEffect(effect)

	if effect.OnApply != nil {
		effect.OnApply(effect)
	}

	return effect
}

func (effect *Effect) Update(dt float64) {
	effect.Entity.Update(dt)
}
