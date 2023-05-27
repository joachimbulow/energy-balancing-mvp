const fs = require("fs");

const { publish } = require("../client");
const {
  getIndex,
  incrementIndex,
  getEnergyApplied,
} = require("./state-redis-client");

const {
  factorInBatteryActions,
  resetBatteryActions,
  checkIfFrequencyIsStabilized,
} = require("../battery-actions");

const FREQUENCY_MEASUREMENT_TOPIC = "frequency_measurements";

const FREQUENCY_PATH = "./measurements/pmu_measurements.json";

/**
 *  This is the nominal static grid data - it will be permanently influenced by batteries
 *  and inertia, but it is the starting point for the simulation.
 */
var frequencyData = [];

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
 * Publishes the current frequency measurements to the broker
 */
async function initializeFrequencyPublication() {
  setInterval(() => {
    publishFrequencyMeasurements();
  }, PUBLISHING_INTERVAL_FREQUENCY_MS);
  async function publishFrequencyMeasurements() {
    const measurements = await getCurrentFrequencyMeasurements();
    console.log("Publishing frequency measurements");

    try {
      publish(FREQUENCY_MEASUREMENT_TOPIC, measurements);
    } catch (error) {
      console.error("Error publishing frequency to Kafka: " + error);
    }
  }
}

/**
 * @returns Data list of the different cities and their frequency measurements at the current inerval counter position
 * Including the effect of battery actions
 */
async function getCurrentFrequencyMeasurements() {
  try {
    const index = (await getIndex()) ?? 0;
    const dataIndex = index * numberOfLocations;

    incrementIndex();

    const currentFrequencyData = frequencyData.slice(
      dataIndex,
      dataIndex + numberOfLocations
    );

    const previouslyAppliedEnergy = await getEnergyApplied();

    const factoredData =
      (await factorInBatteryActions(
        currentFrequencyData,
        previouslyAppliedEnergy
      )) ?? currentFrequencyData;

    await checkIfFrequencyIsStabilized(factoredData);

    const factoredDataWithCurrentTime = factoredData.map((item) => {
      item.timestamp = new Date();
      return item;
    });
    resetBatteryActions();

    return factoredDataWithCurrentTime;
  } catch (error) {
    console.error(
      "Error occured when getting frequency measurements: " + error
    );
  }
}

module.exports = {
  loadFrequencyMeasurements,
  initializeFrequencyPublication,
};
