const fs = require("fs");

const { publish } = require("../client");
const { getIndex, incrementIndex } = require("./state-redis-client");

const {
  factorInBatteryActions,
  resetBatteryActions,
} = require("../battery-actions");

const FREQUENCY_MEASUREMENT_TOPIC = "frequency_measurements";

const FREQUENCY_PATH = "./measurements/pmu_measurements.json";

/**
 *  This is the nominal static grid data - it will be permanently influenced by batteries
 *  and inertia, but it is the starting point for the simulation.
 */
var frequencyData = [];

var currentTime;
var intervalCounter = 0;
const PUBLISHING_INTERVAL_FREQUENCY_MS = 10000; // 10 seconds
var numberOfLocations;

function loadFrequencyMeasurements() {
  fs.readFile(FREQUENCY_PATH, "utf8", (err, data) => {
    if (err) {
      console.error(err);
      return;
    }

    frequencyData = JSON.parse(data);

    numberOfLocations = [
      ...new Set(frequencyData.slice(100).map((item) => item.location)),
    ].length;
  });
}

/**
 * @returns Data list of the different cities and their frequency measurements at the current inerval counter position
 * Including the effect of battery actions
 */
async function getCurrentFrequencyMeasurements() {
  const startIndex = intervalCounter * numberOfLocations;
  const endIndex = startIndex + numberOfLocations;
  intervalCounter++;
  try {
    await incrementIndex();
  } catch (error) {
    console.error("Error incrementing index in Redis");
    console.error(error);
  }

  const currentFrequencyData = frequencyData.slice(startIndex, endIndex);

  const factoredData =
    factorInBatteryActions(currentFrequencyData) ?? currentFrequencyData;

  const factoredDataWithCurrentTime = factoredData.map((item) => {
    item.timestamp = new Date();
    return item;
  });

  resetBatteryActions();

  return factoredDataWithCurrentTime;
}

/**
 * Publishes the current frequency measurements to the broker
 */
async function initializeFrequencyPublication() {
  await loadCheckpoint();

  setInterval(() => {
    publishFrequencyMeasurements();
  }, PUBLISHING_INTERVAL_FREQUENCY_MS);

  async function publishFrequencyMeasurements() {
    const measurements = await getCurrentFrequencyMeasurements();
    console.log("Publishing frequency measurements");

    try {
      publish(FREQUENCY_MEASUREMENT_TOPIC, measurements);
    } catch (error) {
      console.error("Error publishing frequency to Kafka");
      console.error(error);
    }
  }
}

async function loadCheckpoint() {
  var index = 0;

  try {
    index = await getIndex();
  } catch (error) {
    console.error("Error loading checkpoint from Redis");
    console.error(error);
  }
  intervalCounter = index ?? 0;
  console.log(`Checkpoint ${index ? "" : "not"} found, starting from ${index ?? 0}`);
}

module.exports = {
  loadFrequencyMeasurements,
  initializeFrequencyPublication,
};
