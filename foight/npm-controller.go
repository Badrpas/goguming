package foight

import (
	"game/foight/pathfind"
	"game/foight/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

	target_point   cp.Vector
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

		for _, node := range controller.path {
			opts := &ebiten.DrawImageOptions{}
			util.SetDrawOptsColor(opts, color.RGBA{5, 255, 0, 1})
			w, h := _BULLET_IMG.Size()
			opts.GeoM.Translate(float64(w/-2), float64(h/-2))
			opts.GeoM.Scale(0.5, 0.5)
			//size := unit.Game.Nav.GetTileSize() / 2
			//opts.GeoM.Translate(size, size)

			opts.GeoM.Translate(node.X, node.Y)
			unit.Game.TranslateCamera(opts)
			screen.DrawImage(_BULLET_IMG, opts)
		}

		opts := &ebiten.DrawImageOptions{}
		util.SetDrawOptsColor(opts, color.RGBA{255, 0, 0, 1})
		w, h := _BULLET_IMG.Size()
		opts.GeoM.Translate(float64(w/-2), float64(h/-2))
		opts.GeoM.Scale(0.5, 0.5)
		//size := unit.Game.Nav.GetTileSize() / 2
		//opts.GeoM.Translate(size, size)

		opts.GeoM.Translate(controller.target_point.X, controller.target_point.Y)
		unit.Game.TranslateCamera(opts)
		screen.DrawImage(_BULLET_IMG, opts)

		{
			from := TranslatePosFromCamera(unit.Game.Camera, unit.GetPosition())
			to := TranslatePosFromCamera(unit.Game.Camera, controller.target_point)
			//debug.DrawLine(screen, from, to, color.White)
			ebitenutil.DrawLine(screen, from.X, from.Y, to.X, to.Y, color.White)

			dir := cp.Vector{unit.Dx, unit.Dy}.Mult(20)
			to = from.Add(dir)
			ebitenutil.DrawLine(screen, from.X, from.Y, to.X, to.Y, color.RGBA{255, 0, 125, 255})
		}
	}
}

func stateUpdate(unit *Unit) {
	controller := GetNpcController(unit)

	switch controller.state {
	case IDLE:
		unit.Dx, unit.Dy = 0, 0
		// find point to move to
		var closest *Player = nil
		const MAX_DISTANCE = 500.0
		distance := MAX_DISTANCE
		for _, entity := range unit.Game.Entities {
			if player, ok := entity.Holder.(*Player); ok {
				d := unit.GetPosition().Distance(entity.GetPosition())
				if distance > d && isInLos(unit.Entity, entity, unit.Game.Nav) {
					closest = player
					distance = d
				}
			}
		}

		if closest == nil {
			return
		}

		if controller.next_path_calc <= util.TimeNow() {
			path := unit.Game.Nav.FindVec(unit.GetPosition(), closest.GetPosition())
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

			const PATH_CALC_INTERVAL = 1500
			controller.next_path_calc = util.TimeNow() + PATH_CALC_INTERVAL
			controller.target_point = node.Vector
			controller.path = path

			controller.state = MOVE_TO_POINT
		}

		return

	case MOVE_TO_POINT:
		// add impulse in direction
		position := unit.Body.Position()
		diff := controller.target_point.Sub(position)

		if controller.next_path_calc <= util.TimeNow() {
			controller.state = IDLE
		}

		if diff.Length() < 17 {
			if len(controller.path) > 1 {
				controller.path = controller.path[1:]
				controller.target_point = controller.path[0].Vector
				return
			} else {
				controller.state = IDLE
				unit.Dx, unit.Dy = 0, 0
			}

			return
		}

		dir := diff.Normalize()
		unit.Dx, unit.Dy = dir.X, dir.Y
	}

}

func isInLos(unit1, unit2 *Entity, nav *pathfind.Nav) bool {
	size := nav.GetTileSize()
	pos1 := unit1.GetPosition().Mult(1 / size)
	pos2 := unit2.GetPosition().Mult(1 / size)
	line := util.Makeline(int(pos1.X), int(pos1.Y), int(pos2.X), int(pos2.Y))

	for _, point := range line {
		if tile := nav.GetTileAt(point.X, point.Y); tile != nil && tile.Type == pathfind.NavTileWall {
			return false
		}
	}

	return true
}
