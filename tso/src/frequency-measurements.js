const fs = require("fs");
const { publishClient, ensureClientIsConnected } = require("./redis-client");
const { sleep } = require("./utils");
const { factorInBatteryActions } = require("./battery-actions")

const PEM_REQUESTS_TOPIC = "pem_requests";
const PEM_RESPONSES_TOPIC = "pem_responses";
const FREQUENCY_MEASUREMENT_TOPIC = "frequency_measurements";

const FREQUENCY_PATH = "./tso/resources/pmu_measurements.json";
var frequencyData = [];
var currentTime;
var intervalCounter = 0;
const PUBLISHING_INTERVAL_FREQUENCY_MS = 10000; // 10 seconds
var numberOfLocations;

function loadFrequencyMeasurements() {
  console.log("current directory: " + __dirname)
  fs.readFile(FREQUENCY_PATH, "utf8", (err, data) => {
    if (err) {
      console.error(err);
      return;
    }

    frequencyData = JSON.parse(data);

    numberOfLocations = [...new Set(frequencyData.map((item) => item.location))].length;
  });
}

function getCurrentFrequencyMeasurements() {
  currentTime = new Date(frequencyData[intervalCounter * numberOfLocations].timestamp);

  intervalCounter++;

  var currentFrequencyData = frequencyData.filter(
    (measurement) => new Date(measurement.timestamp).getTime() === currentTime.getTime()
  );

  currentFrequencyData.forEach((measurement) => {
    factorInBatteryActions(measurement);
  });

  return currentFrequencyData;
}

function initializeFrequencyPublication() {
  setInterval(() => {
    publishFrequencyMeasurements();
  }, PUBLISHING_INTERVAL_FREQUENCY_MS);

  async function publishFrequencyMeasurements() {
    await ensureClientIsConnected(publishClient);

    const measurements = getCurrentFrequencyMeasurements();

    publishClient.publish(
      FREQUENCY_MEASUREMENT_TOPIC,
      JSON.stringify(measurements, null, 2)
    );
  }
}

module.exports = {
  loadFrequencyMeasurements,
  initializeFrequencyPublication,
};
