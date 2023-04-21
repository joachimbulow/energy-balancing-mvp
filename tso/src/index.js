const { subscribe } = require("./redis-client");

const { loadInertiaMeasurements, initializeInertiaPublication } = require("./inertia-measurements");
const { loadFrequencyMeasurements, initializeFrequencyPublication } = require("./frequency-measurements");
const { handleBatteryAction } = require("./battery-actions");

const PEM_REQUESTS_TOPIC = "pem_requests";
const PEM_RESPONSES_TOPIC = "pem_responses";
const FREQUENCY_MEASUREMENT_TOPIC = "frequency_measurements";
const BATTERY_ACTIONS_TOPIC = "battery_actions";
const INERTIA_MEASUREMENTS_TOPIC = "inertia_measurements";

const brokerTypes = {
  KAFKA: "KAFKA",
  REDIS: "REDIS",
};

initialize();

async function initialize() {
  initializeBroker(brokerTypes.REDIS);

  subscribe(BATTERY_ACTIONS_TOPIC, handleBatteryAction)

  loadFrequencyMeasurements();
  initializeFrequencyPublication();

  loadInertiaMeasurements();
  initializeInertiaPublication();
}