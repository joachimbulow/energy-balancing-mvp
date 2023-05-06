const fs = require('fs');
const path = require('path');
const { randomNormal } = require('d3-random');

const startDate = new Date('2023-01-01T00:00:00Z');
const endDate = new Date('2023-01-01T23:59:59Z');

const SECOND_INTERVAL = 10;

const locations = ['Aarhus', 'Copenhagen', 'Odense', 'Aalborg'];

const voltageMean = 230;
const voltageStdDev = 3;
const currentMean = 10;
const currentStdDev = 2;
const frequencyMean = 50;
const frequencyStdDev = 0.01;
const consumptionMean = 2136.66;
const consumptionStdDev = 547.84;
const productionMean = 2265;
const productionStdDev = 624.16;

const measurements = [];

for (let currentDate = new Date(startDate); currentDate <= endDate; currentDate.setSeconds(currentDate.getSeconds() + SECOND_INTERVAL)) {
  for (let i = 0; i < locations.length; i++) {
    const location = locations[i];
    const measurement = {
      timestamp: currentDate.toISOString(),
      location: location,
      voltage: +(randomNormal(voltageMean, voltageStdDev)()).toFixed(1),
      current: +(randomNormal(currentMean, currentStdDev)()).toFixed(1),
      frequency: +(randomNormal(frequencyMean, frequencyStdDev)()).toFixed(2),
      consumption: +(randomNormal(consumptionMean, consumptionStdDev)()).toFixed(2),
      production: +(randomNormal(productionMean, productionStdDev)()).toFixed(2),
    };
    measurements.push(measurement);
  }
}

const filePath = path.join(__dirname, 'pmu_measurements.json');
fs.writeFileSync(filePath, JSON.stringify(measurements, null, 2), function(err) {
    if(err) {
      console.log(err);
    } else {
      console.log("JSON saved to " + outputFilename);
    }
}); 
