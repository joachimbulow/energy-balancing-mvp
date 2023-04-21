const redis = require("redis");

const REDIS_HOST = "localhost";
const REDIS_PORT = 6379;
const REDIS_CONFIG = {
  socket: {
    host: REDIS_HOST,
    port: REDIS_PORT,
  },
};

console.log("Connecting to Redis at " + REDIS_HOST + ":" + REDIS_PORT);

const publishClient = redis.createClient(REDIS_CONFIG);

publishClient.connect();

publishClient.on("error", (err) => {
  console.log("Error " + err);
});

publishClient.on("connect", () => {
  console.log("Publish client connected to Redis");
});

const subscribeClient = redis.createClient(REDIS_CONFIG);

subscribeClient.connect();

subscribeClient.on("error", (err) => {
  console.log("Error " + err);
});

subscribeClient.on("connect", () => {
  console.log("Subscribe client connected to Redis");
});

async function ensurePublishClientIsConnected() {
  if (publishClient == null || publishClient == undefined) {
    await connectPublishClient(publishClient);
  }
}

async function connectPublishClient(client) {
  client = redis.createClient(REDIS_CONFIG);

  await client.connect();

  client.on("connect", () => {
    console.log("Client connected to Redis");
  });

  client.on("error", (err) => {
    console.log("Error " + err);
  });
}

async function subscribe(topic, handler) {
  subscribeClient.subscribe(topic, (message, topic) => {
    console.log("Received message on channel " + topic + ":\n" + message);
    handler(message);
  });
}

async function publish(topic, message) {
  await ensurePublishClientIsConnected();
  publishClient.publish(
    topic,
    JSON.stringify(message, null, 2)
  );
}

module.exports = {
  subscribe,
  publish,
};
