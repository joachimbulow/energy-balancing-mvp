const { SocketClosedUnexpectedlyError } = require("redis");
var { getCurrentInertia } = require("./measurements/inertia-measurements");

const ACTION = {
  CHARGE: "CHARGE",
  DISCHARGE: "DISCHARGE",
};

const ONE_MILLION = 1000000;
var batteryActions = [];
const energyPacket = 4 / ONE_MILLION; // 4 W in MW
const NOMINAL_FREQUENCY = 50;

// Sum total of all the packets applied
var totalEnergyApplied = 0;

/**
 * Handles the battery actions received from the broker
 * @param {The battery action message} message
 */
function handleBatteryAction(message) {
  batteryActions.push(message);
}

function resetBatteryActions() {
  batteryActions = [];
}

/**
 * Uses Swing equation for calculating how much the battery packets of energy influence the frequency
 * @param { List of measurements obtained from statically generated data } measurement
 * @returns The same list of measurements, but with the frequency adjusted based on the battery actions
 */
function factorInBatteryActions(measurements) {
  if (getCurrentInertia() == 0) {
    console.log("No intertia registered, skipping battery actions");
    return;
  }

  var energyApplied = 0;

  for (const action of batteryActions) {
    if (action.actionType === ACTION.CHARGE) {
      energyApplied -= energyPacket;
    } else if (action.actionType === ACTION.DISCHARGE) {
      energyApplied += energyPacket;
    }
  }
  if (energyApplied != 0) {
    console.log(`Energy change in grid since last refresh: ${energyApplied}`);
  }
  // Update the global state, and use for calculation of new frequency
  totalEnergyApplied += energyApplied;

  console.log(
    `Total change to apply including previous actions: ${totalEnergyApplied}`
  );

  for (const measurement of measurements) {
    var batteryAdjustedFrequency = calculateNewFrequency(
      totalEnergyApplied,
      NOMINAL_FREQUENCY,
      getCurrentInertia(),
      measurement.frequency
    );
    console.log("New frequency = " + batteryAdjustedFrequency);

    measurement.frequency = batteryAdjustedFrequency;
  }

  return measurements;
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
  resetBatteryActions,
};
