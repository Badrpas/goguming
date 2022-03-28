package foight

import (
	"game/foight/util"
	"github.com/jakecoffman/cp"
	"time"
)

type Weapon struct {
	CoolDown   int64
	LastFireAt int64

	Damage int

	Holder *Entity

	CanFireFn func() bool
	EmitFn    func(weapon *Weapon, dir cp.Vector)
}

func (w *Weapon) FireInDirection(dir cp.Vector) {
	if w.IsReady() {
		if w.EmitFn != nil {
			w.EmitFn(w, dir)
		}
	}
}

func (w *Weapon) IsReady() bool {
	can_fire := w.CoolDown+w.LastFireAt <= util.TimeNow()
	if w.CanFireFn != nil {
		can_fire = w.CanFireFn()
	}
	return can_fire
}

func NewWeaponDefault(Holder *Entity) *Weapon {
	return &Weapon{
		CoolDown:   300,
		LastFireAt: util.TimeNow(),
		EmitFn:     DefaultWeaponEmit,

		Holder: Holder,
	}
}

func DefaultWeaponEmit(w *Weapon, dir cp.Vector) {
	e := w.Holder
	unit, is_unit := e.Holder.(*Unit)
	player, is_player := e.Holder.(*Player)
	if is_player {
		unit, is_unit = player.Unit, true
	}

	w.EmitFn = func(w *Weapon, dir cp.Vector) {
		w.LastFireAt = time.Now().UnixMilli()

		b := NewBullet(e.Game, e.X, e.Y)
		b.Team = w.Holder.Team
		b.Shape.Filter.Group = e.Shape.Filter.Group
		b.DrawOpts.ColorM = e.DrawOpts.ColorM
		b.Lifespan = 800
		b.Issuer = e

		b.on_dmg_dealt = func(b *Bullet, to *Entity) {
			if is_player {
				player.AttacksConnectedCount++
			}
			b.Entity.RemoveFromGame()
		}

		e.Game.AddEntity(b.Entity)

		dir = dir.Normalize()
		modifier := 1.
		if is_unit {
			modifier = unit.ForceModifier
		}

		b.Body.SetVelocityVector(
			dir.Mult(modifier * 600.).
				Add(e.Body.Velocity().Mult(1.5)),
		)
	}

	w.EmitFn(w, dir)
}
