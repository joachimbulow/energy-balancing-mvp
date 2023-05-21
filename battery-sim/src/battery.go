package src

import (
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/joachimbulow/pem-energy-balance/src/client"
	"github.com/joachimbulow/pem-energy-balance/src/util"
)

const (
	PEM_REQUESTS_TOPIC          = "pem_requests"
	PEM_RESPONSES_TOPIC         = "pem_responses"
	FREQUENCY_MEASUREMENT_TOPIC = "frequency_measurements"
	BATTERY_ACTIONS_TOPIC       = "battery_actions"
	INERTIA_MEASUREMENT         = "inertia_measurements"
)

type PEMRequest struct {
	ID          string `json:"id"`
	BatteryID   string `json:"batteryId"`
	RequestType string `json:"requestType"`
}

type PEMResponse struct {
	ID           string `json:"id"`
	BatteryID    string `json:"batteryId"`
	ResponseType string `json:"responseType"`
}

type BatteryAction struct {
	ID         string `json:"id"`
	BatteryID  string `json:"batteryId"`
	ActionType string `json:"actionType"`
}

const (
	CHARGE    = "CHARGE"
	DISCHARGE = "DISCHARGE"
)

const (
	GRANTED = "GRANTED"
	DENIED  = "DENIED"
)

const (
	REDIS = "REDIS"
	KAFKA = "KAFKA"
)

const (
	SOC_MEAN           = 0.7
	SOC_STD            = 0.05
	SIGNIFICANT_DIGITS = 4

	BATTERY_CAPACITY_KWH = 13.5
)

var (
	upperBoundBatteryCapacity = util.GetUpperBoundBatteryCapacity()
	lowerBoundBatteryCapacity = util.GetLowerBoundBatteryCapacity()
	requestInterval           = util.GetRequestInterval()
	packetPowerW              = util.GetPacketPowerW()
	packetTimeS               = util.GetPacketTimeS()
	packetEnergyJ             = float64(packetPowerW) * float64(packetTimeS)
	packetKwh                 = float64(packetEnergyJ / 3600000.00)
)

/// -------------------------------------

type Battery struct {
	id                   string
	requestChannel       chan PEMRequest
	responseChannel      chan PEMResponse
	batteryActionChannel chan BatteryAction
	soc                  float64
	latestRequest        PEMRequest
	logger               util.Logger
	busy                 bool
}

func NewBattery(id string, requestChannel chan PEMRequest, responseChannel chan PEMResponse, batteryActionChannel chan BatteryAction) {
	battery := Battery{}
	battery.id = id
	battery.latestRequest = PEMRequest{}
	battery.logger = util.NewLogger(battery.id)
	battery.requestChannel = requestChannel
	battery.responseChannel = responseChannel
	battery.batteryActionChannel = batteryActionChannel

	// Start go routines
	go battery.publishPEMrequests()
	go battery.listenForPEMresponses()
	go battery.logger.Info("Battery started with id/consumer group: %s\n", battery.id)
}

func (battery *Battery) setupClient() client.Client {
	var clientInstance, err = client.NewClient()
	if err != nil {
		battery.logger.Fatalf(err, "Client instance could not be created")
	}
	return clientInstance
}

func GenerateUuid() string {
	return uuid.New().String()
}

// Pem requests
func (battery *Battery) publishPEMrequests() {
	for {
		time.Sleep(requestInterval)
		for battery.busy {
			time.Sleep(time.Second)
		}
		request := battery.getPEMRequest()

		battery.latestRequest = request
		battery.logger.Info("Publishing %s request with id %s and battery id %s\n", request.RequestType, request.ID, request.BatteryID)
		battery.requestChannel <- request
	}
}

func (battery *Battery) getPEMRequest() PEMRequest {
	stateOfCharge := battery.measureSoC()

	var request PEMRequest
	if stateOfCharge < lowerBoundBatteryCapacity {
		request = battery.newRequest(CHARGE)
	} else if stateOfCharge >= lowerBoundBatteryCapacity && stateOfCharge <= upperBoundBatteryCapacity {
		request = battery.probabilisticallyCalculateRequest()
	} else {
		request = battery.newRequest(DISCHARGE)
	}
	return request
}

func (battery *Battery) measureSoC() float64 {
	if battery.soc == 0 {
		// first measurement should be normally distributed around soc mean and std dev
		battery.soc = rand.NormFloat64()*SOC_STD + SOC_MEAN
	}
	battery.logger.Info("Battery state of charge: %.2f\n", battery.soc)
	return battery.soc
}

func (battery *Battery) newRequest(requestType string) PEMRequest {
	return PEMRequest{
		ID:          GenerateUuid(),
		BatteryID:   battery.id,
		RequestType: requestType,
	}
}

// The closer the battery is to the lower bound, the higher the probability of charging.
func (battery *Battery) probabilisticallyCalculateRequest() PEMRequest {
	// Calculate distance to lower and upper bounds
	lowerBoundDistance := battery.soc - lowerBoundBatteryCapacity
	upperBoundDistance := upperBoundBatteryCapacity - battery.soc

	// Calculate probability of charging based on distance to bound
	chargeProbability := math.Abs(lowerBoundDistance) / (math.Abs(lowerBoundDistance) + math.Abs(upperBoundDistance))

	if rand.Float64() < chargeProbability {
		return battery.newRequest(CHARGE)
	} else {
		return battery.newRequest(DISCHARGE)
	}
}

// Pem responses

func (battery *Battery) listenForPEMresponses() {
	for response := range battery.responseChannel {
		if response.ID != battery.latestRequest.ID {
			battery.logger.Info("Received response with id %s, but latest request id is %s. Ignoring response.\n", response.ID, battery.latestRequest.ID)
			continue
		}

		battery.logger.Info("Received %s response with id %s and battery id %s\n", response.ResponseType, response.ID, response.BatteryID)

		if response.ResponseType == GRANTED {
			battery.actOnGrantedRequest(response)
		}
	}
}

func (battery *Battery) actOnGrantedRequest(response PEMResponse) {
	for battery.busy {
		time.Sleep(1 * time.Second)
	}

	request := battery.latestRequest

	if request.RequestType == CHARGE {
		battery.chargePacket()
	} else if request.RequestType == DISCHARGE {
		battery.dischargePacket()
	}
	battery.publishBatteryAction(request.RequestType)
}

func (battery *Battery) chargePacket() {
	battery.logger.Info("Charging packet of + %0.4f kWh.\n", packetKwh)
	battery.updateBattery(packetKwh)
}

func (battery *Battery) dischargePacket() {
	battery.logger.Info("Discharging packet of - %0.4f kWh.\n", packetKwh)
	battery.updateBattery(-packetKwh)
}

func (battery *Battery) updateBattery(chargeAmount float64) {
	battery.busy = true
	currentBatteryCharge := battery.soc * BATTERY_CAPACITY_KWH
	currentBatteryCharge += chargeAmount
	battery.soc = (currentBatteryCharge / BATTERY_CAPACITY_KWH)

	time.Sleep(time.Duration(packetTimeS) * time.Second) // simulate charging/discharging

	battery.logger.Info("After the update the new SoC is: %.4f\n", battery.soc)
	battery.busy = false
}

// Battery actions

func (battery *Battery) publishBatteryAction(actionType string) {
	battery.logger.Info("Publishing battery action: %s\n", actionType)
	action := BatteryAction{
		ID:         battery.id,
		BatteryID:  battery.id,
		ActionType: actionType,
	}
	battery.batteryActionChannel <- action

}
