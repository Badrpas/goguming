## goguming
Couch multiplayer with your mobile phone as controller

### Download
Get prebuilt version as zip archive from [latest release](https://github.com/Badrpas/goguming/releases/latest)

### Build and run
To run locally
```shell
go run .
```

On bind on a custom address

```shell
go run . -addr 0.0.0.0:7331
```

### Custom levels
You can use [Tiled](https://www.mapeditor.org/) to open `.tmx` files in `levels` dir to edit or create new levels

To run a specific level
```shell
game.exe -level levels/mylevel.tmx
```
