const { subscribe } = require("./client");

const {
  loadInertiaMeasurements,
  initializeInertiaPublication,
} = require("./inertia-measurements");
const {
  loadFrequencyMeasurements,
  initializeFrequencyPublication,
} = require("./frequency-measurements");
const { handleBatteryAction } = require("./battery-actions");

const BATTERY_ACTIONS_TOPIC = "battery_actions";

initialize();

async function initialize() {
  subscribe(BATTERY_ACTIONS_TOPIC, handleBatteryAction);

  loadFrequencyMeasurements();
  initializeFrequencyPublication();

  loadInertiaMeasurements();
  initializeInertiaPublication();
}
