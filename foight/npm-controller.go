package foight

import (
	"fmt"
	"game/foight/pathfind"
	"game/foight/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/jakecoffman/cp"
	"image/color"
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

	target_entity  *Entity
	target_pos     cp.Vector
	desired_pos    cp.Vector
	next_path_calc int64

	path []pathfind.NavPathNode
}

func GetNpcController(unit *Unit) *NpcController {
	return NpcControllers[unit]
}
func createNpcController(unit *Unit) *NpcController {
	controller := &NpcController{
		next_path_calc: util.TimeNow(),
	}
	NpcControllers[unit] = controller
	return controller
}
func removeNpcController(unit *Unit) {
	delete(NpcControllers, unit)
}

var DrawNpcNavDebug string

func AddNpcController(unit *Unit) {
	controller := createNpcController(unit)

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

		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			if DrawNpcNavDebug == "true" {
				DrawNpcNavDebug = "false"
			} else {
				DrawNpcNavDebug = "true"
			}
		}
	}

	render := unit.RenderFn
	unit.RenderFn = func(e *Entity, screen *ebiten.Image) {
		if render != nil {
			render(e, screen)
		} else {
			e.Render(screen)
		}

		if DrawNpcNavDebug != "true" {
			return
		}

		draw_node := func(node *pathfind.NavTile, c color.Color) {
			opts := &ebiten.DrawImageOptions{}
			util.SetDrawOptsColor(opts, c)
			w, h := _BULLET_IMG.Size()
			opts.GeoM.Translate(float64(w/-2), float64(h/-2))
			opts.GeoM.Scale(0.5, 0.5)
			size := unit.Game.Nav.GetTileSize() / 2
			opts.GeoM.Translate(size, size)

			opts.GeoM.Translate(float64(node.X*16)-8, float64(node.Y*16)-8)
			unit.Game.TranslateCamera(opts)
			screen.DrawImage(_BULLET_IMG, opts)
		}

		for _, node := range unit.Game.Nav.IterateTiles(pathfind.NavTileFilledGap) {
			draw_node(node, color.RGBA{0, 255, 180, 1})
		}
		for _, node := range unit.Game.Nav.IterateTiles(pathfind.NavTileWall) {
			draw_node(node, color.RGBA{0, 55, 225, 1})
		}
		for _, node := range unit.Game.Nav.IterateTiles(pathfind.NavTileEmpty) {
			draw_node(node, color.RGBA{155, 105, 25, 1})
		}

		draw_marker := func(pos cp.Vector, clr color.Color) {
			opts := &ebiten.DrawImageOptions{}
			util.SetDrawOptsColor(opts, clr)
			w, h := _BULLET_IMG.Size()
			opts.GeoM.Translate(float64(w/-2), float64(h/-2))
			opts.GeoM.Scale(0.5, 0.5)

			opts.GeoM.Translate(pos.X, pos.Y)
			unit.Game.TranslateCamera(opts)
			screen.DrawImage(_BULLET_IMG, opts)
		}

		for _, node := range controller.path {
			draw_marker(node.Vector, color.RGBA{5, 255, 0, 1})
		}

		draw_marker(controller.target_pos, color.RGBA{255, 0, 0, 255})

		draw_marker(controller.desired_pos, color.RGBA{255, 0, 255, 255})

		{
			from := TranslatePosFromCamera(unit.Game.Camera, unit.GetPosition())
			to := TranslatePosFromCamera(unit.Game.Camera, controller.target_pos)
			//debug.DrawLine(screen, from, to, color.White)
			ebitenutil.DrawLine(screen, from.X, from.Y, to.X, to.Y, color.White)

			dir := cp.Vector{unit.Dx, unit.Dy}.Mult(20)
			to = from.Add(dir)
			ebitenutil.DrawLine(screen, from.X, from.Y, to.X, to.Y, color.RGBA{255, 0, 125, 255})
		}

		kda := fmt.Sprintf("%d", controller.state)
		lk := float64(len(kda))
		f := mplusNormalFont
		opts := &ebiten.DrawImageOptions{}
		util.SetDrawOptsColor(opts, color.RGBA{255, 0, 255, 1})

		opts.GeoM.Translate((unit.X - lk*4), (unit.Y + 46))
		unit.Game.TranslateCamera(opts)
		text.DrawWithOptions(screen, kda, f, opts)
	}

}

func stateUpdate(unit *Unit) {
	controller := GetNpcController(unit)

	if controller.target_entity != nil {
		if isInLos(unit.Entity, controller.target_entity, unit.Game.Nav) {
			controller.desired_pos = controller.target_entity.GetPosition()

			diff := controller.target_entity.GetPosition().Sub(unit.GetPosition())
			dir := diff.Normalize()
			unit.Tx, unit.Ty = dir.X, dir.Y

			updateAttack(unit, controller)
		}
	}

	switch controller.state {
	case IDLE:
		unit.Dx, unit.Dy = 0, 0
		// find point to move to
		var closest *Unit = nil
		const MAX_DISTANCE = 500.0
		distance := MAX_DISTANCE
		for _, entity := range unit.Game.Entities {
			if entity.Team == unit.Team {
				continue
			}

			if target_unit, ok := GetUnitFromEntity(entity); ok {
				d := unit.GetPosition().Distance(entity.GetPosition())
				if distance > d && isInLos(unit.Entity, entity, unit.Game.Nav) {
					closest = target_unit
					distance = d
				}
			}
		}

		if closest == nil {
			return
		}

		controller.target_entity = closest.Entity

		if controller.next_path_calc <= util.TimeNow() {
			recalculateUnitPath(unit, controller, closest.GetPosition())
			controller.state = MOVE_TO_POINT
		}

		return

	case MOVE_TO_POINT:
		// add impulse in direction
		position := unit.GetPosition()
		diff := controller.target_pos.Sub(position)

		if controller.next_path_calc <= util.TimeNow() {
			recalculateUnitPath(unit, controller, controller.desired_pos)
		}

		if diff.Length() < 17 {
			if len(controller.path) > 1 {
				controller.path = controller.path[1:]
				controller.target_pos = controller.path[0].Vector
			} else if controller.desired_pos.Distance(position) < 17 {
				controller.state = IDLE
				unit.Dx, unit.Dy = 0, 0
				return
			}
		}

		dir := diff.Normalize()
		unit.Dx, unit.Dy = dir.X, dir.Y
	}

}

func recalculateUnitPath(unit *Unit, controller *NpcController, target cp.Vector) {
	path := unit.Game.Nav.FindVec(unit.GetPosition(), target)
	if len(path) == 0 {
		return
	}

	node := path[0]
	if node.Vector.Distance(unit.GetPosition()) < 30 {
		if len(path) > 1 {
			node = path[1]
			path = path[1:]
		} else {
			return
		}
	}

	const PATH_CALC_INTERVAL = 500
	controller.next_path_calc = util.TimeNow() + PATH_CALC_INTERVAL
	controller.target_pos = node.Vector
	controller.path = path
}

func updateAttack(unit *Unit, controller *NpcController) {
	if controller.target_entity == nil {
		return
	}

	_, ok := controller.target_entity.Holder.(*Unit)
	if !ok {
		_, ok = controller.target_entity.Holder.(*Player)
		if !ok {
			return
		}
	}

	if unit.Weapon == nil || !unit.Weapon.IsReady() {
		return
	}

	unit.Fire()
}
