const redis = require("redis");
const fs = require("fs");
const { constants } = require("buffer");
const sleep = require("./utils").sleep;

const PEM_REQUESTS_TOPIC = "pem_requests";
const PEM_RESPONSES_TOPIC = "pem_responses";
const FREQUENCY_MEASUREMENT_TOPIC = "frequency_measurements";
const BATTERY_ACTIONS_TOPIC = "battery_actions";
const INERTIA_MEASUREMENTS_TOPIC = "inertia_measurements";

const REDIS_HOST = "redis";
const REDIS_PORT = 6379;
const REDIS_CONFIG = {
  socket: {
    host: REDIS_HOST,
    port: REDIS_PORT,
  },
};

var publishClient;

const FREQUENCY_PATH = "../resources/pmu_measurements.json";
var frequencyData = [];
var currentTime;
var intervalCounter = 0;
const PUBLISHING_INTERVAL_FREQUENCY_MS = 10000; // 10 seconds
var numberOfLocations;

const INERTIA_PATH = "../resources/InertiaNordicSyncharea-January-2023.json";
const PUBLISHING_INTERVAL_INERTIA_MS = 10000; // 10 seconds
const ONE_HOUR_MS = 3600000; // 1 hour in milli
var publishingIntervalCounter = 0;
var inertiaData = [];
var inertiaCounter = 0;
var currentInertiaDK2 = 0;

const ACTION = {
  CHARGE: "CHARGE",
  DISCHARGE: "DISCHARGE",
};

const ONE_BILLION = 1000000000;
var batteryActions = [];
const energyPacket = 4 / ONE_BILLION; // 4 W in MW
const NOMINAL_FREQUENCY = 50;

////////////////////////////// MAIN //////////////////////////////

initialize();

//////////////////////// INITIALIZE //////////////////////////////
async function initialize() {
  publishClient = redis.createClient(REDIS_CONFIG);
  await publishClient.connect();

  publishClient.on("error", (err) => {
    console.log("Error " + err);
  });

  const subscribeClient = redis.createClient(REDIS_CONFIG);

  await subscribeClient.connect();

  subscribeClient.subscribe(BATTERY_ACTIONS_TOPIC, (message, topic) => {
    console.log("Received message on channel " + topic + ":\n" + message);
    handleBatteryAction(message);
  });

  loadFrequencyMeasurements();
  initializeFrequencyPublication();

  loadInertiaMeasurements();
  initializeInertiaPublication();
}

////////////////////////////// FREQUENCY //////////////////////////////

function loadFrequencyMeasurements() {
  fs.readFile(FREQUENCY_PATH, "utf8", (err, data) => {
    if (err) {
      console.error(err);
      return;
    }

    frequencyData = JSON.parse(data);

    numberOfLocations = [...new Set(frequencyData.map((item) => item.location))]
      .length;
  });
}

function getCurrentFrequencyMeasurements() {
  currentTime = new Date(
    frequencyData[intervalCounter * numberOfLocations].timestamp
  );

  intervalCounter++;

  var currentFrequencyData = frequencyData.filter(
    (measurement) =>
      new Date(measurement.timestamp).getTime() === currentTime.getTime()
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

////////////////////////////// INERTIA //////////////////////////////

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
  await ensureClientIsConnected(publishClient);

  const measurements = getCurrentInertiaMeasurements();

  console.table(measurements);

  publishClient.publish(
    INERTIA_MEASUREMENTS_TOPIC,
    JSON.stringify(measurements, null, 2)
  );
}

////////////////////////////// BATTERY ACTIONS //////////////////////////////

function handleBatteryAction(message) {
  const action = JSON.parse(message);
  console.log(Date.now() + ": Received action: " + action.actionType);
  batteryActions.push(action);
}

function factorInBatteryActions(measurement) {
  if (currentInertiaDK2 == 0 || batteryActions.length == 0) {
    return;
  }

  var currentFrequency = measurement.frequency;

  var energyBatteriesAddedToGrid = 0;

  batteryActions.forEach((action) => {
    if (action.actionType === ACTION.CHARGE) {
      energyBatteriesAddedToGrid -= energyPacket;
      console.log(
        `Energy charged to battery, and removed from grid ${energyPacket}`
      );
    } else if (action.actionType === ACTION.DISCHARGE) {
      energyBatteriesAddedToGrid += energyPacket;
      console.log(
        `Energy discharged from battery, and removed from grid ${energyPacket}`
      );
    }
  });

  console.log(`Energy change in grid: ${energyBatteriesAddedToGrid}`);

  var frequency = calculateNewFrequency(
    energyBatteriesAddedToGrid,
    NOMINAL_FREQUENCY,
    currentInertiaDK2,
    currentFrequency
  );
  console.log("New frequency = " + frequency);

  // TODO map measurement object. change frequency property to new frequency
  measurement.frequency = frequency;
}

// BASED ON SWING EQUATION
function calculateFrequencyDeviation(
  addedEnergy, // in MW
  nominalFrequency, // in Hz
  inertia // in seconds per megawatt (s/MW)
) {
  console.log(
    `Calculate frequency deviation: added energy = ${addedEnergy}, nominal frequency = ${nominalFrequency}, inertia = ${inertia}, currentInertia = ${currentInertiaDK2}`
  );
  if (nominalFrequency <= 0 || inertia <= 0) {
    throw new Error("Nominal frequency and inertia must be positive numbers.");
  }

  return addedEnergy / (2 * Math.PI * nominalFrequency * inertia);
}

// BASED ON SWING EQUATION
function calculateNewFrequency(
  addedEnergy, // in MW
  nominalFrequency, // in Hz
  inertia, // in seconds per megawatt (s/MW)
  previousFrequency // in Hz
) {
  console.log(
    `Calculate new frequency: added energy = ${addedEnergy}, nominal frequency = ${nominalFrequency}, inertia = ${inertia}, previous frequency = ${previousFrequency}`
  );
  if (nominalFrequency <= 0 || inertia <= 0) {
    throw new Error("Nominal frequency and inertia must be positive numbers.");
  }

  var deviation = addedEnergy / (2 * Math.PI * nominalFrequency * inertia);
  console.log("New deviation = " + deviation);

  // The reason why we subtract the deviation from the previous frequency is that
  // if the added energy is negative (i.e., the batteries are charging), this will cause a decrease in the frequency,
  // while if the added energy is positive (i.e., the batteries are discharging), this will cause an increase in the frequency.
  var newFrequency = previousFrequency - deviation;
  console.log("New frequency = " + newFrequency);

  return newFrequency;
}

function calculateNeededEnergy(
  nominalSystemFrequency, // in  Hz
  frequencyDeviation, // in Hz
  inertia, // seconds per megawatt (s/MW)
  powerImbalance // in MW
) {
  if (nominalSystemFrequency <= 0 || inertia <= 0) {
    throw new Error("Nominal frequency and inertia must be positive numbers.");
  }

  return (
    (frequencyDeviation * inertia * powerImbalance) /
    (2 * Math.PI * nominalSystemFrequency)
  );
}

function testingFormulas() {
  const powerImbalance = 4; // MW
  const addedEnergy = 0.0022918311805232927; // MW
  const nominalFrequency = 50; // Hz
  const frequencyDeviation = 0.02; //Hz
  const inertia = 9; // s/MW

  const energyNeeded = calculateNeededEnergy(
    nominalFrequency,
    frequencyDeviation,
    inertia,
    powerImbalance
  );
  console.log(
    `Needed energy to adjust frequency: ${energyNeeded.toPrecision(4)} MWh`
  );

  const newFrequencyDeviation = calculateFrequencyDeviation(
    addedEnergy,
    nominalFrequency,
    inertia
  );

  console.log(
    `New frequency deviation: ${newFrequencyDeviation.toPrecision(4)} Hz`
  );
}
////////////////////////////// REDIS COMMUNICATION //////////////////////////////

async function connectClient(client) {
  client = redis.createClient(REDIS_CONFIG);

  await client.connect();

  client.on("connect", () => {
    console.log("Publish client connected to Redis");
  });

  client.on("error", (err) => {
    console.log("Error " + err);
  });
}

async function ensureClientIsConnected(client) {
  if (client == null || client == undefined || client.connected == false) {
    await connectClient(client);
  }
}
