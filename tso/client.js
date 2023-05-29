const broker = process.env.BROKER || "KAFKA";
const brokerTypes = {
  KAFKA: "KAFKA",
  REDIS: "REDIS",
  LOGGER: "LOGGER",
};

if (broker.toUpperCase() === brokerTypes.KAFKA) {
  console.log("Using Kafka");
  module.exports = require("./kafka-client");
} else if (broker.toUpperCase() === brokerTypes.REDIS){
  console.log("Using Redis");
  module.exports = require("./redis-client");
} else {
  console.log("Using logger");
  module.exports = require("./logger-client");
}