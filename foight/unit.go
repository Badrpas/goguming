package foight

import (
  "game/foight/util"
  "github.com/hajimehoshi/ebiten/v2"
  "github.com/hajimehoshi/ebiten/v2/text"
  "github.com/jakecoffman/cp"
  "math"
  "strings"
)

type Unit struct {
  *Entity

  Name string

  HP            int32
  is_invincible bool

  Dx, Dy float64
  Tx, Ty float64

  Speed         float64
  ForceModifier float64

  CoolDown       int64
  last_fire_time int64

  Effects []*Effect

  stunned_until int64

  onDeathFn func(self *Unit)
}

func NewUnit(name string, x, y float64, img *ebiten.Image) *Unit {
  unit := &Unit{
    Name: name,
    HP:   DEFAULT_HP,

    Entity: NewEntity(
      x,
      y,
      nil,
      nil,
      img,
    ),

    Speed:         1500,
    ForceModifier: 1,

    CoolDown:       300,
    last_fire_time: util.TimeNow(),
  }

  unit.PreUpdateFn = func(e *Entity, dt float64) {
    unit.UpdateEffects(dt)
    unit.UpdateMovement(dt)
    unit.Entity.PreUpdate(dt)
  }
  unit.UpdateFn = func(e *Entity, dt float64) {
    unit.Update(dt)
  }

  unit.RenderFn = func(e *Entity, screen *ebiten.Image) {
    e.Render(screen)

    info_hp := strings.Repeat("■", int(math.Max(float64(unit.HP), 0)))
    l := float64(len(unit.Name))
    ll := float64(len(info_hp))

    f := mplusNormalFont
    text.Draw(screen, unit.Name, f, int(unit.X-l*4), int(unit.Y-62), unit.color)
    text.Draw(screen, info_hp, f, int(unit.X-ll*2.5), int(unit.Y-42), unit.color)
  }

  unit.OnDmgReceived = func(from *Entity, dmg int32) {
    if unit.is_invincible {
      return
    }

    unit.HP -= dmg

    if unit.HP <= 0 {
      if unit.onDeathFn != nil {
        unit.onDeathFn(unit)
      }

      bullet, ok := from.Holder.(*Bullet)
      if !ok {
        return
      }

      player, ok := bullet.Issuer.Holder.(*Player)
      if !ok {
        return
      }

      player.KillCount++
    }
  }

  unit.onDeathFn = OnDeathDefault

  unit.Entity.Holder = unit

  return unit
}

func (u *Unit) Init(g *Game) {
  body, _ := AddCirclePhysicsToEntity(g, u.Entity)
  body.SetVelocityUpdateFunc(UnitVelocityUpdateFn)

  g.AddEntity(u.Entity)
}

func UnitVelocityUpdateFn(body *cp.Body, gravity cp.Vector, damping float64, dt float64) {
  cp.BodyUpdateVelocity(body, gravity, damping*0.9, dt)
}

func (u *Unit) UpdateMovement(dt float64) {
  if u.IsStunned() {
    return
  }

  mod := dt * u.Speed //* 0.2

  tx := u.Dx * mod
  ty := u.Dy * mod

  impulse := cp.Vector{tx, ty}

  if u.Body != nil {
    if impulse.LengthSq() != 0 {
      u.Body.ApplyImpulseAtLocalPoint(impulse, cp.Vector{})
    }
  }
}

func (u *Unit) Update(dt float64) {
  if u.Tx != 0 || u.Ty != 0 {
    u.Angle = math.Atan2(float64(u.Ty), float64(u.Tx)) + math.Pi/2
  } else if u.Dx != 0 || u.Dy != 0 {
    u.Angle = math.Atan2(float64(u.Dy), float64(u.Dx)) + math.Pi/2
  }

  u.Entity.Update(dt)
}

func (u *Unit) SetInvincible(duration int64) {
  u.is_invincible = true
  precolor := u.color

  u.DrawOpts.ColorM.Scale(0.4, 0.4, 0.4, 1)

  iteration := 1
  interval_id := u.TimeManager.SetInterval(func() {
    iteration++
    if iteration%2 == 0 {
      u.DrawOpts.ColorM.Scale(0.4, 0.4, 0.4, 1)
    } else {
      u.SetColor(precolor)
    }
  }, 300)

  u.TimeManager.SetTimeout(func() {
    u.TimeManager.ClearInterval(interval_id)
    u.is_invincible = false
    u.SetColor(precolor)
  }, duration)
}


func (u *Unit) StunFor(ms int64) {
  u.stunned_until = util.TimeNow() + ms
}
func (u *Unit) IsStunned() bool {
  return util.TimeNow() < u.stunned_until
}


func (u *Unit) UpdateEffects(dt float64) {
  count := len(u.Effects)
  row_length := int(math.Round(math.Sqrt(float64(count))))

  for idx, effect := range u.Effects {
    effect.Update(dt)

    x := idx % row_length
    y := idx / row_length
    effect.X, effect.Y = float64(x-row_length/2)*16, float64(y-row_length/2)*16
  }
}

func (u *Unit) AddEffect(effect *Effect) {
  u.Effects = append(u.Effects, effect)

  effect.Parent = u.Entity

  if effect.Game != u.Game {
    u.Game.AddEntity(effect.Entity)
  }
}

func (u *Unit) RemoveEffect(effect *Effect) {
  if effect.OnCease != nil {
    effect.OnCease(effect)
  }

  for idx, x := range u.Effects {
    if x == effect {
      u.Effects = append(u.Effects[:idx], u.Effects[idx+1:]...)
      break
    }
  }

  if effect.Game != nil {
    effect.RemoveFromGame()
  }
  effect.Target = nil
  effect.Data = nil
}

func OnDeathDefault(unit *Unit) {
  unit.RemoveFromGame()
}

