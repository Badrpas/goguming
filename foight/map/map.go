package levelmap

import (
	"game/foight"
	"github.com/jakecoffman/cp"
	"github.com/lafriks/go-tiled"
	"github.com/solarlune/ldtkgo"
	"log"
)

func LoadToGameLdtk(path string, game *foight.Game) error {
	file, err := ldtkgo.Open(path)
	if err != nil {
		log.Println(err)
		return err
	}

	grid := file.Levels[0].Layers[0].IntGrid

	for _, cell := range grid {
		x, y := cell.Position[0], cell.Position[1]
		b := foight.NewBlock(float64(x), float64(y))
		b.Init(game)
	}

	log.Println(file)

	return nil
}

func LoadToGameTiled(path string, game *foight.Game) error {
	gameMap, err := tiled.LoadFile(path)
	if err != nil {
		log.Println(err)
		return err
	}

	layer := gameMap.Layers[0]
	tiles := layer.Tiles

	for idx, tile := range tiles {
		if tile.Nil {
			continue
		}

		x := 16. * float64(idx%(gameMap.Width))
		y := 16. * float64(idx/(gameMap.Width))

		b := foight.NewBlock(x, y)
		b.Init(game)
	}

	for _, objectGroup := range gameMap.ObjectGroups {
		var points = make([]cp.Vector, len(objectGroup.Objects))
		switch objectGroup.Name {
		case "player_spawn_points":
			game.PlayerSpawnPoints = points
		case "item_spawn_points":
			game.ItemSpawnPoints = points
		default:
			log.Println("Unknown object group name", objectGroup.Name)
			continue
		}

		for i, info := range objectGroup.Objects {
			points[i] = cp.Vector{info.X, info.Y}
		}
	}

	return nil
}
