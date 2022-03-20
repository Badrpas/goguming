package foight

import (
	"flag"
	"game/foight/net"
	"github.com/gorilla/websocket"
	"image/color"
	"log"
	"net/http"
	_ "net/http/pprof"
	"regexp"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

type Message struct {
	from    string
	payload string
}

func shakeHand(c *websocket.Conn) (name string, clr color.Color, err error) {
	{
		log.Println("Reading NAME")
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Print("WS name error:", err)
			return "", color.White, err
		}
		name = string(message)
		log.Println(name)
	}

	{
		log.Println("Reading COLOR")
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Print("WS color error:", err)
			return "", color.White, err
		}
		clr = color.RGBA{message[0], message[1], message[2], 255}
		log.Println(message[0], message[1], message[2], message[3])
	}

	return name, clr, nil
}

func getWsHandler(game *Game) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade error:", err)
			return
		}

		defer c.Close()

		name, color, err := shakeHand(c)
		if err != nil {
			log.Println("Error from handshake:", err)
			return
		}

		log.Printf("New player connected: [%s] color: %X", name, color)

		player := (<-game.QueueJob(func() interface{} {
			return NewPlayer(game, name, color)
		})).(*Player)

		for {
			_, message, err := c.ReadMessage()

			if err != nil {
				log.Println("read err:", err)

				game.QueueJobVoid(func() {
					game.RemoveEntity(player.Entity)
				})
				break
			}

			player.messages <- net.DecodeUpdateMessage(message)
		}

	}

}

func getRootHandler() func(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("./static/"))
	return func(w http.ResponseWriter, r *http.Request) {
		if ok, err := regexp.MatchString("\\.js$", r.URL.Path); ok && err == nil {
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		}
		fs.ServeHTTP(w, r)
	}
}

func RunApi(game *Game) {
	http.HandleFunc("/ws", getWsHandler(game))
	http.HandleFunc("/", getRootHandler())

	log.Println("Starting HTTP server on", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
