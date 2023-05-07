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
const frequencyStdDev = 0.3;
const consumptionMean = 2136.66;
const consumptionStdDev = 547.84;
const productionMean = 2265;
const productionStdDev = 624.16;

var measurements = [];
var prevValues = {};

var degreeOfChange = 1;
var frequencyOfChange = 0;

function resetExperiment() {
  measurements = [];
  prevValues = {};

  degreeOfChange = 1;
  frequencyOfChange = 0;
}

function randomWalk(value, stdDev) {
  if (Math.random() < frequencyOfChange) {
    return value + (Math.random() * 2 - 1) * stdDev + stdDev * degreeOfChange;
  }
  return value + (Math.random() * 2 - 1) * stdDev;
}

// TODO add timestamp to filename -${new Date().toISOString()}
function createExperiment(fileName = `pmu-measurements.json`) {
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

  fileName = "data/" + fileName;
  let filePath = path.join(__dirname, fileName);
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
  resetExperiment();
}

createExperiment();

degreeOfChange = 10;
frequencyOfChange = 0.001;

createExperiment("sudden-change.json");