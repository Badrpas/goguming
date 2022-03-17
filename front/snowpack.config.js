
/** @type {import("snowpack").SnowpackUserConfig } */
const config = {
  mount: {
    static: { url: '/', static: true },
    src: { url: '/' },
  },
  plugins: [
  ],
  routes: [
    /* Enable an SPA Fallback in development: */
    // {"match": "routes", "src": ".*", "dest": "/index.html"},
  ],
  optimize: {
    /* Example: Bundle your final build: */
    "bundle": true,
  },
  packageOptions: {
    /* ... */
  },
  devOptions: {
    /* ... */
    // hostname: '192.168.50.53',
  },
  buildOptions: {
    out: '../static'
  },
  env: {
  }

};

module.exports = config;
