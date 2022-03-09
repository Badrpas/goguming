import { setupController } from './controller.js';

const getUsername = async () => {
  return 'yoba';
}

const getColor = async () => {
  return 0x123456;
}

window.addEventListener("load",  async function(evt) {
  const output = document.getElementById("output");

  var print = function(message) {
    const d = document.createElement("div");
    d.textContent = message;
    output.appendChild(d);
    output.scroll(0, output.scrollHeight);
  };

  const username = await getUsername();
  const color = await getColor();

  const ws = new WebSocket(`ws://${window.location.host}/ws`);
  ws.onopen = function () {
    print("OPEN");
    ws.send(username);
    ws.send(new Uint32Array([color]));

    setupController(ws);
  }

  ws.onclose = function(evt) {
    print("CLOSE");
  }
  ws.onmessage = function(evt) {
    print("RESPONSE: " + evt.data);
  }
  ws.onerror = function(evt) {
    print("ERROR: " + evt.data);
  }


});
