import nipplejs from './lib/nipplejs/index.js';

const Length = 4;

export const setupController = (ws) => {
  const arr = new Uint8Array(Length);
  const sendUpdate = () => {
    ws.send(arr)
  }


  const joystickL = nipplejs.create({
    zone: document.getElementById('left'),
    mode: 'dynamic',
    position: { left: '37%', top: '50%' },
    color: 'blue',
    size: 200
  });

  const joystickR = nipplejs.create({
    zone: document.getElementById('right'),
    mode: 'dynamic',
    position: { left: '63%', top: '50%' },
    color: 'blue',
    size: 200
  });

  [joystickL, joystickR].forEach((joy, i) => {
    const d = i*2;

    joy.on('move', (_, info) => {
      const { x, y } = info.vector;
      arr[d  ] = Math.round(x * 50 + 50);
      arr[d+1] = Math.round(y *-50 + 50);
      sendUpdate();
    });

    joy.on('end', () => {
      arr[d  ] = 50;
      arr[d+1] = 50;
      sendUpdate();
    })
  });

}


