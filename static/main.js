
const getUsername = async () => {
  return 'yoba';
}

window.addEventListener("load",  async function(evt) {
  const output = document.getElementById("output");
  let ws;

  var print = function(message) {
    const d = document.createElement("div");
    d.textContent = message;
    output.appendChild(d);
    output.scroll(0, output.scrollHeight);
  };

  const username = await getUsername();

  ws = new WebSocket("ws://localhost:8080/ws");
  ws.onopen = function(evt) {
    print("OPEN");
    ws.send(username);

    const Length = 2;
    const arr = new Uint8Array(Length);
    const sendUpdate = (x, y) => {
      arr[0] = 1+x;
      arr[1] = 1+y;
      ws.send(arr)
    }

    document.addEventListener('keydown', e => {
      switch (e.key) {
        case 'ArrowUp':
          sendUpdate(0, -1)
          break;
        case 'ArrowDown':
          sendUpdate(0, +1)
          break;
        case 'ArrowLeft':
          sendUpdate(-1, 0)
          break;
        case 'ArrowRight':
          sendUpdate(+1, 0)
          break;
      }
    });
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


});
