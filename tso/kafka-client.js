const { Kafka } = require("kafkajs");

const broker_url = process.env.BROKER_URL || "localhost";
const broker_port = process.env.BROKER_PORT || "9092";

const KAFKA_SERVER = broker_url + ":" + broker_port;
const GROUP_ID = "TSO";

let kafka;
let producer;
let consumer;

connect();

async function connect() {
  kafka = new Kafka({ brokers: [KAFKA_SERVER] });
  producer = kafka.producer();
  consumer = kafka.consumer({ groupId: GROUP_ID });
  await producer.connect();
  await consumer.connect();
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
  await consumer.subscribe({ topic });
  await consumer.run({
    eachMessage: async ({ message }) => {
      if (message.topic === topic) {
        handler(message.value.toString());
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
