package foight

import (
	"github.com/gorilla/websocket"
	"image/color"

	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
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
				<-game.QueueJob(func() interface{} {
					game.RemoveEntity(player.Entity)
					return 0
				})
				log.Println("read err:", err)
				break
			}

			player.messages <- decodeUpdateMessage(message)
		}

	}

}

func RunApi(game *Game) {
	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/ws", getWsHandler(game))
	http.Handle("/", http.FileServer(http.Dir("./static/")))

	log.Println("Starting HTTP server on", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
