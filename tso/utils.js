const fs = require("fs");

function sleep(ms) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
}

// create a function to log using the other log function. It should append the current datetime to the csvData
// and then call the log function.


function log(data = [], headerFields = [], logFilePath = `${new Date().toISOString().slice(0, 10)}-log.txt`) {
  var csvHeader = "dateTime;";
  headerFields.forEach((element) => {
    csvHeader += `${element};`;
  });
  csvHeader += "\n";

  var csvData = `${new Date()};`;
  data.forEach((element) => {
    csvData += `${element};`;
  });

  fs.appendFileSync(logFilePath, `${csvHeader}${csvData}\n`);
}

function executeAndTimeRequest(request) {
  console.time("responseTime");
  request();
  console.timeEnd("responseTime");
  fs.appendFileSync(
    "log.txt",
    `Response time: ${console.timeEnd("responseTime")} ms\n`
  );
}

// TODO IDEA: create a function that will read the log file and return the last line.
// This will be used to get the last frequency, so that the program can continue from where it left off.

module.exports = {
  sleep,
  log,
  executeAndTimeRequest,
};
