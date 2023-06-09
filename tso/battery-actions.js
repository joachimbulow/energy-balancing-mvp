const { SocketClosedUnexpectedlyError } = require("redis");
var { getCurrentInertia } = require("./measurements/inertia-measurements");
const {
  resetEnergyApplied,
  incrementEnergyAppliedBy,
  setEnergyApplied,
} = require("./measurements/state-redis-client");

const ACTION = {
  CHARGE: "CHARGE",
  DISCHARGE: "DISCHARGE",
};

const ONE_BILLION = 1000000000; // Wattseconds = joule so 1000000000 joules = 1 GigaWattsecond

var batteryActions = [];

const PACKET_TIME_S = parseInt(process.env.PACKET_TIME_S || 5 * 60); // Default to 5 minutes
const PACKET_POWER_W = parseInt(process.env.PACKET_POWER_W || 4000); // Default to 4000 watts

const ENERGY_PACKET_J = PACKET_POWER_W * PACKET_TIME_S; // Joules / wattseconds

const ENERGY_PACKET_GIGA_WATT_SECONDS = ENERGY_PACKET_J / ONE_BILLION;

const NOMINAL_FREQUENCY = 50;

var previousMeasurements = [];
/**
 * Handles the battery actions received from the broker
 * @param {The battery action message} message
 */
function handleBatteryAction(message) {
  batteryActions.push(message);
}

function resetBatteryActions() {
  console.log(
    `Resetting ${batteryActions.length} battery actions after publish`
  );
  batteryActions = [];
}

/**
 * Uses Swing equation for calculating how much the battery packets of energy influence the frequency
 * @param { List of measurements obtained from statically generated data } measurement
 * @returns The same list of measurements, but with the frequency adjusted based on the battery actions
 */
async function factorInBatteryActions(measurements, previouslyAppliedEnergy) {
  if (getCurrentInertia() == 0) {
    console.log("No intertia registered, skipping battery actions");
    return;
  }

  var energyApplied = previouslyAppliedEnergy;

  for (const action of batteryActions) {
    if (action.actionType === ACTION.CHARGE) {
      energyApplied -= ENERGY_PACKET_GIGA_WATT_SECONDS;
    } else if (action.actionType === ACTION.DISCHARGE) {
      energyApplied += ENERGY_PACKET_GIGA_WATT_SECONDS;
    }
  }

  // Add value to state in redis
  await setEnergyApplied(energyApplied);

  console.log(
    `Total change to apply including previous actions since last stabilization: ${energyApplied} GWs`
  );

  for (const measurement of measurements) {
    var batteryAdjustedFrequency = calculateNewFrequency(
      energyApplied,
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
 * GWs = GigaWatt-seconds
 * @param {The amount of energy added in GWs} addedEnergy
 * @param {System nominal frequency in Hz} nominalFrequency
 * @param {System current inertia in GWs} inertia // I think this is Gigawatt seconds (GWs) given https://energinet.dk/media/4xinfk4y/ffr-ig-justification-report.pdf
 * @param {The current system frequency before applying battery action in Hz} previousFrequency
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

  console.log(`Applied deviation = ${appliedDeviation} Hz`);

  // appliedDeviation (Δf) is negative when energy is added to the system
  // Therefore, we subtract the deviation from the previous frequency
  var newFrequency = previousFrequency - appliedDeviation;

  return newFrequency;
}

async function checkIfFrequencyIsStabilized(newMeasurements) {
  if (previousMeasurements == null || previousMeasurements.length == 0) {
    console.log(
      `No previous measurements, skipping frequency stabilization check`
    );
    previousMeasurements = newMeasurements;
    return;
  }

  for (var i = 0; i < previousMeasurements.length; i++) {
    var previousFrequency = previousMeasurements[i].frequency;
    var newFrequency = newMeasurements[i].frequency;
    if (
      (previousFrequency > NOMINAL_FREQUENCY &&
        newFrequency <= NOMINAL_FREQUENCY) ||
      (previousFrequency < NOMINAL_FREQUENCY &&
        newFrequency >= NOMINAL_FREQUENCY)
    ) {
      console.log(
        `Frequency has stabilized, resetting total energy applied to 0`
      );
      await resetEnergyApplied();
      break;
    }
  }
  previousMeasurements = newMeasurements;
}

module.exports = {
  handleBatteryAction,
  factorInBatteryActions,
  resetBatteryActions,
  checkIfFrequencyIsStabilized,
};
