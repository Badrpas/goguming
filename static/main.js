import { setupController } from './controller.js';

const HsbToRgb = (h, s, b) => {
  s /= 100;
  b /= 100;
  const k = (n) => (n + h / 60) % 6;
  const f = (n) => b * (1 - s * Math.max(0, Math.min(k(n), 4 - k(n), 1)));
  return [255 * f(5), 255 * f(3), 255 * f(1)].map(Math.round);
};

const getUsername = async () => {
  return 'yoba';
}

const getColor = async () => {
  const color = HsbToRgb(Math.round(Math.random() * 360), 100, 100)

  console.log(color);
  const lr = (color[0] << 16) + (color[1] << 8) + color[2];
  console.log(lr.toString(16));
  return lr

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
