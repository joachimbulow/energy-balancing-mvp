const fs = require("fs");
const { publish } = require("./redis-client");

const INERTIA_MEASUREMENTS_TOPIC = "inertia_measurements";

const INERTIA_PATH = "./tso/resources/InertiaNordicSyncharea-January-2023.json";
const PUBLISHING_INTERVAL_INERTIA_MS = 10000; // 10 seconds
const ONE_HOUR_MS = 3600000; // 1 hour in milli
var publishingIntervalCounter = 0;
var inertiaData = [];
var inertiaCounter = 0;
var currentInertiaDK2 = 0;

function loadInertiaMeasurements() {
  fs.readFile(INERTIA_PATH, "utf8", (err, data) => {
    if (err) {
      console.error(err);
      return;
    }

    inertiaData = JSON.parse(data);

    inertiaData.sort((a, b) => {
      return new Date(a.HourUTC) - new Date(b.HourUTC);
    });
  });
}

// for one hour we publish the same measurement so we don't increment the counter
// every hour we publish a new measurement and we increment the counter by 1
function getCurrentInertiaMeasurements() {
  publishingIntervalCounter++;
  if (inertiaCounter === inertiaData.length - 1) {
    inertiaCounter = 0;
  }
  // if one hour has passed we increment the counter
  var timePassedMs = publishingIntervalCounter * PUBLISHING_INTERVAL_INERTIA_MS;
  if (ONE_HOUR_MS === timePassedMs) {
    console.log(
      "1 hour passed. Incrementing inertia counter: " + inertiaCounter
    );
    inertiaCounter++;
    publishingIntervalCounter = 0;
  }
  currentInertiaDK2 = inertiaData[inertiaCounter].InertiaDK2GWs;
  return inertiaData[inertiaCounter];
}

function initializeInertiaPublication() {
  setInterval(() => {
    publishInertiaMeasurements();
  }, PUBLISHING_INTERVAL_INERTIA_MS);
}

async function publishInertiaMeasurements() {
  const measurements = getCurrentInertiaMeasurements();
  console.table(measurements);

  publish(INERTIA_MEASUREMENTS_TOPIC, measurements);
}

module.exports = {
  loadInertiaMeasurements,
  initializeInertiaPublication,
  currentInertiaDK2: currentInertiaDK2,
};
