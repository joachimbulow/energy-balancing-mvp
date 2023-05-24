const { Kafka } = require("kafkajs");

const broker_url = process.env.BROKER_URL || "localhost";
const broker_port = process.env.BROKER_PORT || "29092";

const KAFKA_SERVER = broker_url + ":" + broker_port;
const GROUP_ID = "TSO";

let kafka;
let producer;
let consumer1;
let consumer2;

connect();

// Utilzing two consumers in the same group
async function connect() {
  kafka = new Kafka({ brokers: [KAFKA_SERVER] });
  producer = kafka.producer();
  consumer1 = kafka.consumer({ groupId: GROUP_ID });
  consumer2 = kafka.consumer({ groupId: GROUP_ID });
  await producer.connect();
  await consumer1.connect();
  await consumer2.connect();
}

async function ensurePublishClientIsConnected() {
  if (!producer) {
    await connectPublishClient();
  }
}

async function connectPublishClient() {
  producer = kafka.producer();
  await producer.connect();
  console.log("Producer connected to Kafka");
}

async function subscribe(topic, handler) {
  // Run one consumer
  await consumer1.subscribe({ topic });
  consumer1.run({
    eachMessage: async (message) => {
      const msg = JSON.parse(message.message.value);
      if (!msg.heartbeat) {
        handler(msg);
      }
    },
  });
  await consumer2.subscribe({ topic });
  consumer2.run({
    eachMessage: async (message) => {
      const msg = JSON.parse(message.message.value);
      if (!msg.heartbeat) {
        handler(msg);
      }
    },
  });
}

async function publish(topic, message) {
  await ensurePublishClientIsConnected();
  await producer.send({
    topic,
    messages: [{ value: JSON.stringify(message, null, 2) }],
  });
}

module.exports = {
  subscribe,
  publish,
};
