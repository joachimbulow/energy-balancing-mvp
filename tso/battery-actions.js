var { currentInertiaDK2 } = require("./inertia-measurements");

const ACTION = {
  CHARGE: "CHARGE",
  DISCHARGE: "DISCHARGE",
};

const ONE_BILLION = 1000000000;
var batteryActions = [];
const energyPacket = 4 / ONE_BILLION; // 4 W in MW
const NOMINAL_FREQUENCY = 50;

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

  measurement.frequency = frequency;
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

module.exports = {
  handleBatteryAction,
  factorInBatteryActions,
};