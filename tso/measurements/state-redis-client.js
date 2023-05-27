const redis = require("redis");

const REDIS_HOST = process.env.REDIS_BROKER || "localhost";
const REDIS_PORT = process.env.REDIS_PORT || 6379;
const REDIS_CONFIG = {
  socket: {
    host: REDIS_HOST,
    port: REDIS_PORT,
  },
};

console.log("Connecting to Redis at " + REDIS_HOST + ":" + REDIS_PORT);

const client = redis.createClient(REDIS_CONFIG);
var connected = false;

async function connect() {
  if (connected) return;
  connected = true;
  await client.connect();

  client.on("connect", () => {
    console.log("Client connected to Redis");
  });

  client.on("error", (err) => {
    console.log("Redis error " + err);
  });
}

async function getIndex() {
  await connect();
  return await client.GET("index");
}

async function incrementIndex() {
  await connect();
  return await client.INCR("index");
}

async function getEnergyApplied() {
  await connect();
  const energy = await client.GET("energy_applied");
  return parseFloat(energy ?? 0);
}

async function setEnergyApplied(energy) {
  await connect();
  return await client.SET("energy_applied", energy.toString());
}

async function resetIndex() {
  await connect();
  return await client.SET("index", 0);
}

async function resetEnergyApplied() {
  await connect();
  return await client.SET("energy_applied", "0");
}

module.exports = {
  incrementIndex,
  getIndex,
  getEnergyApplied,
  setEnergyApplied,
  resetEnergyApplied,
};
