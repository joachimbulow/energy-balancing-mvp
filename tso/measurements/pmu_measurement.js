const fs = require("fs");
const path = require("path");
const { randomNormal } = require("d3-random");

const startDate = new Date("2023-01-01T00:00:00Z");
const endDate = new Date("2023-01-01T23:59:59Z");

const SECOND_INTERVAL = 10;

const locations = ["Aarhus", "Copenhagen", "Odense", "Aalborg"];

const voltageMean = 230;
const voltageStdDev = 3;
const currentMean = 10;
const currentStdDev = 2;
const frequencyMean = 50;
const frequencyStdDev = 0.3;
const consumptionMean = 2136.66;
const consumptionStdDev = 547.84;
const productionMean = 2265;
const productionStdDev = 624.16;

const measurements = [];
const prevValues = {};

function randomWalk(value, stdDev) {
  return value + (Math.random() * 2 - 1) * stdDev;
}

for (
  let currentDate = new Date(startDate);
  currentDate <= endDate;
  currentDate.setSeconds(currentDate.getSeconds() + SECOND_INTERVAL)
) {
  for (let i = 0; i < locations.length; i++) {
    const location = locations[i];

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
        voltageStdDev / 10
      ).toFixed(2),
      current: +randomWalk(
        prevValues[location].current,
        currentStdDev / 10
      ).toFixed(1),
      frequency: +randomWalk(
        prevValues[location].frequency,
        frequencyStdDev / 10
      ).toFixed(5),
      consumption: +randomWalk(
        prevValues[location].consumption,
        consumptionStdDev / 10
      ).toFixed(2),
      production: +randomWalk(
        prevValues[location].production,
        productionStdDev / 10
      ).toFixed(2),
    };

    prevValues[location] = measurement;

    measurements.push(measurement);
  }
}

const filePath = path.join(__dirname, "pmu_measurements.json");
fs.writeFileSync(
  filePath,
  JSON.stringify(measurements, null, 2),
  function (err) {
    if (err) {
      console.log(err);
    } else {
      console.log("JSON saved to " + outputFilename);
    }
  }
);
