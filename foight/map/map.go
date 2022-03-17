package levelmap

import (
	"game/foight"
	"github.com/jakecoffman/cp"
	"github.com/lafriks/go-tiled"
	"image/color"
	"log"
	"math/rand"
)

func LoadToGameTiled(path string, game *foight.Game) error {
	gameMap, err := tiled.LoadFile(path)
	if err != nil {
		log.Println(err)
		return err
	}

	layer := gameMap.Layers[0]
	tiles := layer.Tiles

	cell_w := float64(gameMap.TileWidth)
	cell_h := float64(gameMap.TileHeight)

	for idx, tile := range tiles {
		if tile.Nil {
			continue
		}

		x := cell_w * float64(idx%(gameMap.Width))
		y := cell_h * float64(idx/(gameMap.Width))

		b := foight.NewBlock(x, y)
		b.Init(game)

		c := uint8(155 + rand.Int()%100)
		b.SetColor(color.RGBA{c, c, c, 255})
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
