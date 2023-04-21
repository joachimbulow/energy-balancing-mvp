const fs = require('fs');
const path = require('path');
const { randomNormal } = require('d3-random');

// Define the start and end date of the dataset
const startDate = new Date('2023-01-01T00:00:00Z');
const endDate = new Date('2023-01-01T23:59:59Z');

const SECOND_INTERVAL = 10;

// Define the list of locations where measurements will be taken
const locations = ['Aarhus', 'Copenhagen', 'Odense', 'Aalborg'];

// Define the mean and standard deviation for each measurement type
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

// Define an empty array to store the measurements
const measurements = [];

// Loop through each second between the start and end date
for (let currentDate = new Date(startDate); currentDate <= endDate; currentDate.setSeconds(currentDate.getSeconds() + SECOND_INTERVAL)) {
  // Loop through each location and generate a measurement
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
    // Append the measurement to the array of measurements
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
