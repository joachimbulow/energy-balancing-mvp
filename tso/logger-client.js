async function subscribe(topic, handler) {
    console.log("LOGGER | SUBSCRIBE topic: ", topic, " | handler: ", handler)
}

async function publish(topic, message) {
    console.log("LOGGER | PUBLISH topic: ", topic, " | message: ", JSON.stringify(message, null, 2))
}

module.exports = {
  subscribe,
  publish,
};
