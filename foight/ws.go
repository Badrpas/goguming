package foight

import (
  "github.com/gorilla/websocket"

  "flag"
  "log"
  "net/http"
)


var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

type Message struct {
  from string
  payload string
}

//var Messages = make(chan Message, 1024);



func getWsHandler(game *Game) func(w http.ResponseWriter, r *http.Request) {
  return func (w http.ResponseWriter, r *http.Request) {

    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
      log.Print("upgrade error:", err)
      return
    }

    defer c.Close()

    var name string;
    var player *Player = nil;
    {
      _, message, err := c.ReadMessage()
      if err != nil {
        log.Print("WS name error:", err)
        return;
      }
      name = string(message)
      log.Printf("New player connected: [%s]", name)

      player = game.AddPlayer(name);
    }


    for {
      mt, message, err := c.ReadMessage()
      if err != nil {
        log.Println("read:", err)
        break
      }

      player.messages <- decodeUpdateMessage(message)

      log.Printf("recv: %s", message)
      err = c.WriteMessage(mt, message)
      if err != nil {
        log.Println("write:", err)
        break
      }
    }

  }

}

func RunApi(game *Game) {
  flag.Parse()
  log.SetFlags(0)

  http.HandleFunc("/ws", getWsHandler(game))
  http.Handle("/", http.FileServer(http.Dir("./static/")))



  log.Println("Starting HTTP server")
  log.Fatal(http.ListenAndServe(*addr, nil))
}


