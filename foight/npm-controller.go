package foight

import (
  "github.com/jakecoffman/cp"
  "math/rand"
)

type NpcControllerStorage map[*Unit]*NpcController
var NpcControllers NpcControllerStorage

func init() {
  NpcControllers = NpcControllerStorage{}
}

type NpcControllerState uint
const (
  IDLE NpcControllerState = iota
  MOVE_TO_POINT
)

type NpcController struct {
  state NpcControllerState

  target_point cp.Vector
}


func GetNpcController(unit *Unit) *NpcController {
  return NpcControllers[unit]
}
func createNpcController(unit *Unit) *NpcController {
  controller := &NpcController{
  }
  NpcControllers[unit] = controller
  return controller
}
func removeNpcController(unit *Unit) {
  delete(NpcControllers, unit)
}


func AddNpcController(unit *Unit) {
  createNpcController(unit)

  on_remove := unit.OnRemove
  unit.OnRemove = func(entity *Entity) {
    if on_remove != nil {
      on_remove(entity)
    }
    removeNpcController(unit)
  }

  pre_update := unit.PreUpdateFn
  unit.PreUpdateFn = func(e *Entity, dt float64) {
    stateUpdate(unit)

    if pre_update != nil {
      pre_update(e, dt)
    } else {
      unit.PreUpdate(dt)
    }
  }
}

func stateUpdate(unit *Unit) {
  controller := GetNpcController(unit)

  switch controller.state {
  case IDLE:
    // find point to move to
    controller.target_point = cp.Vector{rand.Float64() * 1600, rand.Float64() * 1000}
    controller.state = MOVE_TO_POINT
    return

  case MOVE_TO_POINT:
    // add impulse in direction
    position := unit.Body.Position()
    diff := controller.target_point.Sub(position)
    if diff.Length() < 10 {
      controller.state = IDLE
      return
    }

    dir := diff.Normalize()
    unit.Dx, unit.Dy = dir.X, dir.Y
  }

}
