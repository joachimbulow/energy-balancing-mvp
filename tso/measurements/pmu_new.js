const math = require("mathjs");
const fs = require("fs");
const path = require("path");


// Parameters
const theta = 0.1;
// timestep
const dt = 10;
// simulation time
const T = 100000;
// number of timesteps
const N = Math.floor(T / dt);

// mu = mean
// sigma = standard deviation

// Voltage
const mu_v = 230;
const sigma_v = 3;
var v = new Array(N).fill(0);
v[0] = mu_v;

// Current
const mu_i = 9;
const sigma_i = 1;
var i = new Array(N).fill(0);
i[0] = mu_i;

// Frequency
const mu_f = 50;
const sigma_f = 0.01;
var f = new Array(N).fill(0);
f[0] = mu_f;

// Consumption
const mu_c = 2136.66;
const sigma_c = 547.84;
var c = new Array(N).fill(0);
c[0] = mu_c;

// Production
const mu_p = 2265;
const sigma_p = 624.16;
var p = new Array(N).fill(0);
p[0] = mu_p;

const locations = ["Aarhus", "Copenhagen", "Odense", "Aalborg"];

var measurements = [];

// Generate values for each timestep
for (let t = 1; t < N; t++) {
  for (let j = 0; j < locations.length; j++) {
    const location = locations[j];

    // Voltage
    const dv = theta * (mu_v - v[t - 1]) * dt + sigma_v * math.sqrt(dt) * math.random();
    v[t] = v[t - 1] + dv;

    // Current
    const di = theta * (mu_i - i[t - 1]) * dt + sigma_i * math.sqrt(dt) * math.random();
    i[t] = i[t - 1] + di;

    // Frequency
    const df = theta * (mu_f - f[t - 1]) * dt + sigma_f * math.sqrt(dt) * math.random();
    f[t] = f[t - 1] + df;

    // Consumption
    const dc = theta * (mu_c - c[t - 1]) * dt + sigma_c * math.sqrt(dt) * math.random();
    c[t] = c[t - 1] + dc;

    // Production
    const dp = theta * (mu_p - p[t - 1]) * dt + sigma_p * math.sqrt(dt) * math.random();
    p[t] = p[t - 1] + dp;

    const measurement = {
      timestamp: t,
      location: location,
      voltage: parseFloat(v[t].toFixed(2)),
      current: parseFloat(i[t].toFixed(2)),
      frequency: parseFloat(f[t].toFixed(4)),
      consumption: parseFloat(c[t].toFixed(0)),
      production: parseFloat(p[t].toFixed(0)),
    };
    measurements.push(measurement);
  }
}

const filePath = path.join(__dirname, '/pmu_new.json');
let jsonOut = "[" + measurements.map(measurement => JSON.stringify(measurement)).join(",") + "]";
fs.writeFileSync(filePath, jsonOut, function(err) {
    if(err) {
      console.log(err);
    } else {
      console.log("JSON saved to " + outputFilename);
    }
});

console.log("JSON saved to " + filePath);
