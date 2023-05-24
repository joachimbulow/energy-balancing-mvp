const fs = require("fs");
const path = require("path");

const startDate = new Date("2023-01-01T00:00:00Z");
const endDate = new Date("2023-01-01T23:59:59Z");

const SECOND_INTERVAL = 10;

const locations = ["Aarhus", "Copenhagen", "Odense", "Aalborg"];

const voltageMean = 230;
const voltageStdDev = 3;
const currentMean = 10;
const currentStdDev = 2;
const frequencyMean = 50;
const frequencyStdDev = 0.038;
const consumptionMean = 2136.66;
const consumptionStdDev = 547.84;
const productionMean = 2265;
const productionStdDev = 624.16;

var measurements = [];
var prevValues = {};

var degreeOfChange = 0;
var frequencyOfChange = 0;
var probabilityOfDecreasing = 0.5;
var counter = 0;
var changeDir = false;
var changeDirectionAfter = 10 * 10;

const today = new Date().toISOString().slice(0, 16);

function resetExperiment() {
  measurements = [];
  prevValues = {};

  degreeOfChange = 0;
  frequencyOfChange = 0;
  probabilityOfDecreasing = 0.5;
  counter= 0;
}

function randomWalk(value, stdDev, randomness = Math.random(), chanceOfChange = Math.random(), direction = Math.random() < probabilityOfDecreasing ? -1 : 1) {
  let localDeviation = Math.random() * 3;

  if (chanceOfChange < frequencyOfChange) {
    return value + (randomness * 2 - 1) * stdDev + stdDev * degreeOfChange * localDeviation * direction;
  }
  return value + (randomness * 2 - 1) * stdDev * localDeviation * direction;
}

function createExperiment(fileName = `pmu-measurements-${today}.json`) {
  for (
    let currentDate = new Date(startDate);
    currentDate <= endDate;
    currentDate.setSeconds(currentDate.getSeconds() + SECOND_INTERVAL)
  ) {
    if (changeDir == true && counter % changeDirectionAfter == 0) {
      console.log("Changing direction: " + currentDate);
      probabilityOfDecreasing = 1 - probabilityOfDecreasing;
    }

    var randomness = Math.random();
    var chanceOfChange = Math.random();
    let direction = Math.random() < probabilityOfDecreasing ? -1 : 1;

    for (let i = 0; i < locations.length; i++) {
      const location = locations[i];

      createMeasurement(location, currentDate, randomness, chanceOfChange, direction);
    }
    counter++;
  }

  fileName = "data/" + fileName;
  let filePath = path.join(__dirname, fileName);
  fs.writeFileSync(
    filePath,
    JSON.stringify(measurements, null, 2),
    function (err) {
      if (err) {
        console.log(err);
      } else {
        console.log("JSON saved to " + filePath);
      }
    }
  );
  console.log("JSON saved to " + filePath);
  resetExperiment();
}

function createMeasurement(location, currentDate, randomness, chanceOfChange, direction) {
  if (!prevValues[location]) {
    prevValues[location] = {
      voltage: voltageMean,
      current: currentMean,
      frequency: frequencyMean,
      consumption: consumptionMean,
      production: productionMean,
    };
  }

  const measurement = {
    timestamp: currentDate.toISOString(),
    location: location,
    voltage: +randomWalk(
      prevValues[location].voltage,
      voltageStdDev / 10,
      randomness,
      chanceOfChange,
      direction
    ).toFixed(2),
    current: +randomWalk(
      prevValues[location].current,
      currentStdDev / 10,
      randomness,
      chanceOfChange,
      direction
    ).toFixed(1),
    frequency: +randomWalk(
      prevValues[location].frequency,
      frequencyStdDev / 10,
      randomness,
      chanceOfChange,
      direction
    ).toFixed(5),
    consumption: +randomWalk(
      prevValues[location].consumption,
      consumptionStdDev / 10,
      randomness,
      chanceOfChange,
      direction
    ).toFixed(2),
    production: +randomWalk(
      prevValues[location].production,
      productionStdDev / 10,
      randomness,
      chanceOfChange,
      direction
    ).toFixed(2),
  };

  prevValues[location] = measurement;

  measurements.push(measurement);
}

changeDir = true;
degreeOfChange = 8;
frequencyOfChange = 0.11;
probabilityOfDecreasing = 0.9;
createExperiment();