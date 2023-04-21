const { subscribeClient } = require("./redis-client");

var { currentInertiaDK2 } = require("./inertia-measurements");
const { loadInertiaMeasurements, initializeInertiaPublication } = require("./inertia-measurements");
const { loadFrequencyMeasurements, initializeFrequencyPublication } = require("./frequency-measurements");
const { handleBatteryAction } = require("./battery-actions");


const PEM_REQUESTS_TOPIC = "pem_requests";
const PEM_RESPONSES_TOPIC = "pem_responses";
const FREQUENCY_MEASUREMENT_TOPIC = "frequency_measurements";
const BATTERY_ACTIONS_TOPIC = "battery_actions";
const INERTIA_MEASUREMENTS_TOPIC = "inertia_measurements";

initialize();

async function initialize() {
  subscribeClient.subscribe(BATTERY_ACTIONS_TOPIC, (message, topic) => {
    console.log("Received message on channel " + topic + ":\n" + message);
    handleBatteryAction(message);
  });

  loadFrequencyMeasurements();
  initializeFrequencyPublication();

  loadInertiaMeasurements();
  initializeInertiaPublication();
}