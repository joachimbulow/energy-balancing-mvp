var { currentInertiaDK2 } = require("./measurements/inertia-measurements");

const ACTION = {
  CHARGE: "CHARGE",
  DISCHARGE: "DISCHARGE",
};

const ONE_BILLION = 1000000000;
var batteryActions = [];
const energyPacket = 4 / ONE_BILLION; // 4 W in MW
const NOMINAL_FREQUENCY = 50;

// Sum total of all the packets applied
var totalEnergyApplied = 0;

/**
 * Handles the battery actions received from the broker
 * @param {The battery action message} message
 */
function handleBatteryAction(message) {
  const action = JSON.parse(message);
  batteryActions.push(action);
}

/**
 * Uses Swing equation for calculating how much the battery packets of energy influence the frequency
 * @param { List of measurements obtained from statically generated data } measurement
 * @returns The same list of measurements, but with the frequency adjusted based on the battery actions
 */
function factorInBatteryActions(measurement) {
  if (currentInertiaDK2 == 0 || batteryActions.length == 0) {
    return;
  }

  var currentFrequency = measurement.frequency;

  var energyApplied = 0;

  for (const action of batteryActions) {
    if (action.actionType === ACTION.CHARGE) {
      energyApplied -= energyPacket;
    } else if (action.actionType === ACTION.DISCHARGE) {
      energyApplied += energyPacket;
    }
  }

  // Reset battery actions
  batteryActions = [];

  console.log(`Energy change in grid since last refresh: ${energyApplied}`);

  // Update the global state, and use for calculation of new frequency
  totalEnergyApplied += energyApplied;

  console.log(`Total change to apply: ${energyApplied}`);

  var frequency = calculateNewFrequency(
    totalEnergyApplied,
    NOMINAL_FREQUENCY,
    currentInertiaDK2,
    currentFrequency
  );
  console.log("New frequency = " + frequency);

  measurement.frequency = frequency;
}

/**
 *
 * @param {The amount of energy added in MW} addedEnergy
 * @param {System nominal frequency in Hz} nominalFrequency
 * @param {System current inertia in seconds per megawatts (s/MW)} inertia
 * @param {The current system frequency before applying battery action} previousFrequency
 * @returns
 */
function calculateNewFrequency(
  addedEnergy,
  nominalFrequency,
  inertia,
  previousFrequency
) {
  if (nominalFrequency <= 0 || inertia <= 0) {
    throw new Error("Nominal frequency and inertia must be positive numbers.");
  }

  // Use Swing equation
  // AddedEnergy is ΔP, also known as, deviation in power (in this case applied by the batteries)
  var appliedDeviation =
    addedEnergy / (2 * Math.PI * nominalFrequency * inertia);

  // appliedDeviation (Δf) is negative when energy is added to the system
  // Therefore, we subtract the deviation from the previous frequency
  var newFrequency = previousFrequency - appliedDeviation;

  return newFrequency;
}

module.exports = {
  handleBatteryAction,
  factorInBatteryActions,
};
