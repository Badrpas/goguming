package levelmap

import (
	"game/foight"
	"game/foight/pathfind"
	imagestore "game/img"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"github.com/lafriks/go-tiled"
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

func LoadToGameTiled(path string, game *foight.Game) error {
	file, err := tiled.LoadFile(path)
	if err != nil {
		log.Println(err)
		return err
	}

	ProcessMapProperties(game, file.Properties)

	layer := file.Layers[0]
	tiles := layer.Tiles

	game.Nav.SetSize(file.Width, file.Height)
	game.Nav.SetTileSize(float64(file.TileWidth))

	cell_w := float64(file.TileWidth)
	cell_h := float64(file.TileHeight)

	for idx, tile := range tiles {
		if tile.Nil {
			continue
		}

		idx_x := (idx % (file.Width))
		idx_y := (idx / (file.Width))
		x := cell_w * float64(idx_x)
		y := cell_h * float64(idx_y)

		img := GetImage(tile)
		b := foight.NewBlock(x, y, img)
		b.Init(game)

		c := uint8(155 + rand.Int()%100)
		b.SetColor(color.RGBA{c, c, c, 255})

		SetWallAroundPoint(game.Nav, idx_x, idx_y, 1)
	}

	game.Nav.FixHolesWithActorSize(2)
	game.Nav.Init()

	for _, objectGroup := range file.ObjectGroups {
		switch objectGroup.Name {
		case "player_spawn_points":
			game.PlayerSpawnPoints = CollectPositions(objectGroup)
		case "item_spawn_points":
			game.ItemSpawnPoints = CollectPositions(objectGroup)

		case "npc_spawn_info":
			spawn_info := CollectNpcInfos(objectGroup)
			for _, info := range spawn_info {
				game.SpawnNpc(info)
			}

		case "flags":
			ProcessFlagInfos(game, objectGroup)

		default:
			log.Println("Unknown object group name", objectGroup.Name)
			continue
		}

	}

	return nil
}

func ProcessMapProperties(game *foight.Game, properties *tiled.Properties) {
	if properties == nil {
		return
	}
	mode := properties.GetString("game_mode")
	if mode == "coop" {
		game.Mode = foight.GameModeCoop
	}
}

func CollectNpcInfos(group *tiled.ObjectGroup) []*foight.NpcSpawnInfo {
	infos := make([]*foight.NpcSpawnInfo, len(group.Objects))

	for idx, object := range group.Objects {
		info := &foight.NpcSpawnInfo{
			Pos:    cp.Vector{object.X, object.Y},
			Name:   object.Name,
			Weapon: "default",
			Color:  color.RGBA{255, 0, 200, 255},
			HP:     5,
			Team:   -1,
		}

		for _, property := range object.Properties {
			val := property.Value
			switch property.Name {
			case "color":
				r, _ := strconv.ParseUint(val[3:5], 16, 8)
				g, _ := strconv.ParseUint(val[5:7], 16, 8)
				b, _ := strconv.ParseUint(val[7:9], 16, 8)
				info.Color = color.RGBA{uint8(r), uint8(g), uint8(b), 255}
			case "weapon":
				info.Weapon = val
			case "name":
				info.Name = val
			case "hp":
				info.HP, _ = strconv.Atoi(val)
			case "team":
				info.Team, _ = strconv.Atoi(val)
			}
		}

		infos[idx] = info
	}

	return infos
}

func CollectPositions(objectGroup *tiled.ObjectGroup) []cp.Vector {
	points := make([]cp.Vector, len(objectGroup.Objects))

	for i, info := range objectGroup.Objects {
		points[i] = cp.Vector{info.X, info.Y}
	}

	return points
}

func ProcessFlagInfos(game *foight.Game, group *tiled.ObjectGroup) {
	flag_handler := foight.NewFlagHandler()
	game.AddEntity(flag_handler.Entity)

	for _, info := range group.Objects {
		pos := cp.Vector{info.X, info.Y}
		flag := foight.NewFlag(pos, flag_handler)
		flag.Init(game)
	}
}

func SetWallAroundPoint(nav *pathfind.Nav, x, y, radius int) {
	//nav.SetWall(x, y)

	for i := 0; i <= radius; i++ {
		nav.SetWall(x+radius, y+i)
		nav.SetWall(x+radius, y-i)
		nav.SetWall(x-radius, y+i)
		nav.SetWall(x-radius, y-i)

		nav.SetWall(x+i, y+radius)
		nav.SetWall(x-i, y+radius)
		nav.SetWall(x+i, y-radius)
		nav.SetWall(x-i, y-radius)
	}

}

func GetImage(t *tiled.LayerTile) *ebiten.Image {
	for _, prototile := range t.Tileset.Tiles {
		if prototile.ID == t.ID {
			img_name := prototile.Image.Source
			img_name = strings.Replace(img_name, "../img/", "", 1)
			img, ok := imagestore.Images[img_name]
			if !ok {
				print("Can't find image", img_name)
			}
			return img
		}
	}
	return nil
}
