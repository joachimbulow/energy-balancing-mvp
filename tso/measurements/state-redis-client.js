const redis = require("redis");

const REDIS_HOST = process.env.BROKER_URL || "localhost";
const REDIS_PORT = process.env.BROKER_PORT || 6379;
const REDIS_CONFIG = {
  socket: {
    host: REDIS_HOST,
    port: REDIS_PORT,
  },
};

console.log("Connecting to Redis at " + REDIS_HOST + ":" + REDIS_PORT);

const client = redis.createClient(REDIS_CONFIG);

async function getIndex() {
  return await client.GET("index");
}

async function incrementIndex() {
  return await client.INCR("index");
}

async function resetIndex() {
  return await client.SET("index", 0);
}
