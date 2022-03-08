package foight

import (
  "github.com/gorilla/websocket"

  "flag"
  "fmt"
  "html/template"
  "log"
  "net/http"
  "time"
)


var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

type Message struct {
  from string
  payload string
}

var Messages = make(chan Message, 1024);

func home(w http.ResponseWriter, r *http.Request) {
  err := homeTemplate.Execute(w, "ws://"+r.Host+"/ws")

  if err != nil {
    log.Println("Some error while rendering home", err)
    return
  }
}

func (g *Game) readMessages () {
  var msg *Message = nil
  select {
  case _msg := <-Messages:
    msg = &_msg
  default:
  }

  if msg != nil {
    log.Printf("Received message from [%s]: %s\n", msg.from, msg.payload)
    g.last_message = fmt.Sprintf("[%s]: %s", msg.from, msg.payload)
  }
}

func getWsHandler(game *Game) func(w http.ResponseWriter, r *http.Request) {
  return func (w http.ResponseWriter, r *http.Request) {

    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
      log.Print("upgrade error:", err)
      return
    }

    defer c.Close()

    var name string;
    {
      _, message, err := c.ReadMessage()
      if err != nil {
        log.Print("WS name error:", err)
        return;
      }
      name = string(message)
      log.Printf("New player connected: [%s]", name)

      game.AddPlayer(name);
    }


    for {
      mt, message, err := c.ReadMessage()
      if err != nil {
        log.Println("read:", err)
        break
      }

      Messages <- Message {
        from: name,
        payload: string(message),
      }

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
  http.HandleFunc("/", home)

  go func() {
    Messages <- Message{
      from: "autokek",
      payload: "heheboi",
    }

    time.Sleep(2 * time.Second);

    Messages <- Message{
      from: "autokek 2 msg",
      payload: "lul xd",
    }
  }()

  log.Println("Starting HTTP server")
  log.Fatal(http.ListenAndServe(*addr, nil))
}



var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))