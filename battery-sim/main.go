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

	// Setup the common masterClient that all batteries will use, as well as a few slave clients for consumption only
	masterClient, masterError := client.NewClient() // Handles consumption and production
	slaveClient1, slaveError1 := client.NewClient() // Handles only consumption
	slaveClient2, slaveError2 := client.NewClient() // Handles only consumption
	slaveClient3, slaveError3 := client.NewClient() // Handles only consumption
	slaveClient4, slaveError4 := client.NewClient() // Handles only consumption

	if masterError != nil || slaveError1 != nil || slaveError2 != nil || slaveError3 != nil || slaveError4 != nil {
		panic("Failed to initialize a Kafka client")
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

	// Producer
	go publishPEMrequests(requestChannel, masterClient)
	go publishBatteryActions(batteryActionChannel, masterClient)

	// Consumer only clients
	go listenForPEMresponses(responseChannelMap, masterClient)
	go listenForPEMresponses(responseChannelMap, slaveClient1)
	go listenForPEMresponses(responseChannelMap, slaveClient2)
	go listenForPEMresponses(responseChannelMap, slaveClient3)
	go listenForPEMresponses(responseChannelMap, slaveClient4)

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
		go func(req src.PEMRequest) {
			jsonRequest, err := json.Marshal(req)
			if err != nil {
				logger.ErrorWithMsg("Marshaling of pem request message failed", err)
				return
			}
			client.Publish(src.PEM_REQUESTS_TOPIC, req.BatteryID, string(jsonRequest))
		}(request)
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
		go func(act src.BatteryAction) {
			jsonAction, err := json.Marshal(act)
			if err != nil {
				logger.ErrorWithMsg("Marshaling of battery action message failed", err)
				return
			}
			client.Publish(src.BATTERY_ACTIONS_TOPIC, act.BatteryID, string(jsonAction))
		}(action)
	}
}
