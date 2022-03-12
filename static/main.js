import { setupController } from './controller.js';

const getUsername = async () => {
  return 'yoba';
}

const getColor = async () => {
  return 0xFF00FF;
}

window.addEventListener("load",  async function(evt) {
  const output = document.getElementById("output");

  const print = function(message) {
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

  ws.onclose = async function(evt) {
    print("CLOSE");

    while (true) {
      await new Promise(resolve => setTimeout(resolve, 3000))
      console.log('retrying connection');
      try {
        await fetch(window.location.href)
        console.log('connection established. reloading');
        return location.reload()
      } catch (err) {
        console.log(`Couldn't get the root. retrying`);
        await new Promise(resolve => setTimeout(resolve, 3000))
      }
    }
  }
  ws.onmessage = function(evt) {
    print("RESPONSE: " + evt.data);
  }
  ws.onerror = function(evt) {
    print("ERROR: " + evt.data);
  }


});
