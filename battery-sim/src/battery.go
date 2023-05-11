package src

import (
	"encoding/json"
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
	UpperBoundBatteryCapacity = 0.8
	LowerBoundBatteryCapacity = 0.2

	SOC_MEAN           = 0.7
	SOC_STD            = 0.05
	SIGNIFICANT_DIGITS = 4

	BATTERY_CAPACITY_KWH = 13.5

	PacketPowerW    = 4000
	PacketTimeS     = 5 * 60
	PACKET_ENERGY_J = PacketPowerW * PacketTimeS

	PACKET_KWH = PACKET_ENERGY_J / 3600000.00

	SENDING_INTERVAL_NS = 10 * time.Second

	CHARGE_DISCHARGE_INTERVAL_MS = 10000
)

type Battery struct {
	id            string
	client        client.Client
	soc           float64
	latestRequest PEMRequest
	logger        util.Logger
	busy          bool
}

func NewBattery() {
	battery := Battery{}
	battery.id = generateUuid()
	battery.client = battery.setupClient()
	battery.latestRequest = PEMRequest{}
	battery.logger = util.NewLogger(battery.id)
	go battery.client.Listen(PEM_RESPONSES_TOPIC, battery.handlePEMresponse)
	go battery.publishPEMrequests()
	battery.logger.Info("Battery started")
}

func (battery *Battery) setupClient() client.Client {
	var clientInstance, err = client.NewClient()
	if err != nil {
		battery.logger.Fatalf(err, "Client instance could not be created")
	}
	return clientInstance
}

func generateUuid() string {
	return uuid.New().String()
}

func (battery *Battery) handlePEMresponse(params ...[]byte) {
	message := params[1]

	response := PEMResponse{}
	if err := json.Unmarshal(message, &response); err != nil {
		battery.logger.ErrorWithMsg("Unmarshaling of message failed", err)
		return
	}

	if response.ID != battery.latestRequest.ID {
		return
	}

	battery.logger.Info("Received %s response with id %s and battery id %s\n", response.ResponseType, response.ID, response.BatteryID)

	if response.ResponseType == GRANTED {
		battery.actOnGrantedRequest(response)
	} else if response.ResponseType == DENIED {
		battery.logger.Info("Request with id %s denied\n", response.ID)
	}
}

func (battery *Battery) publishPEMrequests() {
	for {
		time.Sleep(SENDING_INTERVAL_NS)
		for battery.busy {
			time.Sleep(time.Second)
		}
		request := battery.getPEMRequest()

		jsonRequest, err := json.Marshal(request)
		if err != nil {
			battery.logger.Fatal(err)
		}

		battery.latestRequest = request
		battery.logger.Info("Sending %s request with id %s and battery id %s\n", request.RequestType, request.ID, request.BatteryID)
		battery.client.Publish(PEM_REQUESTS_TOPIC, request.BatteryID, string(jsonRequest))
	}
}

func (battery *Battery) getPEMRequest() PEMRequest {
	stateOfCharge := battery.measureSoC()

	var request PEMRequest
	if stateOfCharge < LowerBoundBatteryCapacity {
		request = battery.newRequest(CHARGE)
	} else if stateOfCharge >= LowerBoundBatteryCapacity && stateOfCharge <= UpperBoundBatteryCapacity {
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
		ID:          generateUuid(),
		BatteryID:   battery.id,
		RequestType: requestType,
	}
}

// The closer the battery is to the lower bound, the higher the probability of charging.
func (battery *Battery) probabilisticallyCalculateRequest() PEMRequest {
	// Calculate distance to lower and upper bounds
	lowerBoundDistance := battery.soc - LowerBoundBatteryCapacity
	upperBoundDistance := UpperBoundBatteryCapacity - battery.soc

	// Calculate probability of charging based on distance to bound
	chargeProbability := math.Abs(lowerBoundDistance) / (math.Abs(lowerBoundDistance) + math.Abs(upperBoundDistance))

	if rand.Float64() < chargeProbability {
		return battery.newRequest(CHARGE)
	} else {
		return battery.newRequest(DISCHARGE)
	}
}

func (battery *Battery) actOnGrantedRequest(response PEMResponse) {
	battery.logger.Info("Request with id %s approved\n", response.ID)
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
	battery.logger.Info("Charging packet of + %0.4f kWh.\n", PACKET_KWH)
	battery.updateBattery(PACKET_KWH)
}

func (battery *Battery) dischargePacket() {
	battery.logger.Info("Discharging packet of - %0.4f kWh.\n", PACKET_KWH)
	battery.updateBattery(-PACKET_KWH)
}

func (battery *Battery) updateBattery(chargeAmount float64) {
	battery.busy = true
	currentBatteryCharge := battery.soc * BATTERY_CAPACITY_KWH
	currentBatteryCharge += chargeAmount
	battery.soc = (currentBatteryCharge / BATTERY_CAPACITY_KWH)
	time.Sleep(CHARGE_DISCHARGE_INTERVAL_MS * time.Millisecond) // simulate charging/discharging
	battery.logger.Info("After the update the new SoC is: %.4f\n", battery.soc)
	battery.busy = false
}

func (battery *Battery) publishBatteryAction(actionType string) {
	battery.logger.Info("Publishing battery action: %s\n", actionType)
	action := BatteryAction{
		ID:         generateUuid(),
		BatteryID:  battery.id,
		ActionType: actionType,
	}
	json, err := json.Marshal(action)
	if err != nil {
		battery.logger.Fatal(err)
	}
	battery.client.Publish(BATTERY_ACTIONS_TOPIC, battery.id, string(json))
}
