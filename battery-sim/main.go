package main

import (
	"encoding/json"
	"time"

	"github.com/joachimbulow/pem-energy-balance/src"
	"github.com/joachimbulow/pem-energy-balance/src/client"
	"github.com/joachimbulow/pem-energy-balance/src/util"
)

var (
	nBatteries = util.GetNumberOfBatteries()
	logger     = util.Logger{}
)

func main() {
	initializeBatteries()
}

func initializeBatteries() {
	// All request flow "up-stream" from batteries through this channel
	requestChannel := make(chan src.PEMRequest)
	// All responses flow "down-stream" to the right batteries through these channels
	responseChannelMap := make(map[string]chan src.PEMResponse)
	// All battery actions flow "up-stream" from batteries through this channel
	batteryActionChannel := make(chan src.BatteryAction)

	// Setup the common client that all batteries will use
	client, error := client.NewClient()

	if error != nil {
		println("Failed to initialize the main client")
		panic(error)
	}

	for i := 0; i < nBatteries; i++ {
		batteryId := src.GenerateUuid()
		// Create channels for the battery
		responseChannel := make(chan src.PEMResponse)

		// Add the channels to a map, so we can keep track
		responseChannelMap[batteryId] = responseChannel

		// To spread out the spawn of the batteries so they spam a bit less
		time.Sleep(100 * time.Millisecond)
		go startBattery(batteryId, requestChannel, responseChannel, batteryActionChannel)
	}

	// Start the channel communication go routines
	go publishPEMrequests(requestChannel, client)
	go listenForPEMresponses(responseChannelMap, client)
	go publishBatteryActions(batteryActionChannel, client)

	// To keep routines running' we start sleepin'
	for {
		time.Sleep(1 * time.Second)
	}
}

func startBattery(id string, requestChannel chan src.PEMRequest, responseChannel chan src.PEMResponse, batteryActionChannel chan src.BatteryAction) {
	src.NewBattery(id, requestChannel, responseChannel, batteryActionChannel)
}

// Send out PEM requests when the batteries requests it through channels
func publishPEMrequests(requestsChannel chan src.PEMRequest, client client.Client) {
	for request := range requestsChannel {
		jsonRequest, err := json.Marshal(request)
		if err != nil {
			logger.ErrorWithMsg("Marshaling of pem request message failed", err)
			continue
		}
		client.Publish(src.PEM_REQUESTS_TOPIC, request.BatteryID, string(jsonRequest))
	}
}

// Listen for PEM responses
func listenForPEMresponses(responseChannelMap map[string]chan src.PEMResponse, client client.Client) {
	client.Listen(src.PEM_RESPONSES_TOPIC, util.GetConsumerGroupId(), func(params ...[]byte) { handlePemResponse(responseChannelMap, params...) })
}

// Send out the response to the correct channel based on the id
func handlePemResponse(responseChannelMap map[string]chan src.PEMResponse, params ...[]byte) {
	message := params[1]

	response := src.PEMResponse{}
	if err := json.Unmarshal(message, &response); err != nil {
		logger.ErrorWithMsg("Unmarshaling of response message failed", err)
		return
	}

	channel, ok := responseChannelMap[response.BatteryID] // Only the messages meant for our batteries should be published to the batteries
	if !ok {
		return
	}

	channel <- response
}

// Send out battery actions
func publishBatteryActions(batteryActionChannel chan src.BatteryAction, client client.Client) {
	for action := range batteryActionChannel {
		jsonAction, err := json.Marshal(action)
		if err != nil {
			logger.ErrorWithMsg("Marshaling of battery action message failed", err)
			continue
		}
		client.Publish(src.BATTERY_ACTIONS_TOPIC, action.BatteryID, string(jsonAction))
	}
}
