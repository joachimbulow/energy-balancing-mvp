package src

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/joachimbulow/pem-energy-balance/src/broker"
	"github.com/joachimbulow/pem-energy-balance/src/util"
)

const (
	PEM_REQUESTS_TOPIC          = "pem_requests"
	PEM_RESPONSES_TOPIC         = "pem_responses"
	FREQUENCY_MEASUREMENT_TOPIC = "frequency_measurements"
	BATTERY_ACTIONS_TOPIC       = "battery_actions"
	INERTIA_MEASUREMENT         = "inertia_measurements"
)

var (
	logger util.Logger
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

	BATTERY_CAPACITY_KWH = 2

	PacketPowerW    = 4000
	PacketTimeS     = 5 * 60
	PACKET_ENERGY_J = PacketPowerW * PacketTimeS

	PACKET_KWH = float64(PACKET_ENERGY_J / 3600000)

	SENDING_INTERVAL_NS = 10 * time.Second

	CHARGE_DISCHARGE_INTERVAL_MS = 10000
)

type Battery struct {
	id       string
	broker   broker.Broker
	soc      float64
	requests map[string]PEMRequest
	logger   util.Logger
}

var (
	busy = false
)

func NewBattery() Battery {
	battery := Battery{}
	battery.id = generateUuid()
	battery.broker = battery.setupBroker()
	battery.requests = make(map[string]PEMRequest)
	battery.logger = util.NewLogger(battery.id)
	go battery.broker.Listen(PEM_RESPONSES_TOPIC, battery.handlePEMresponse)
	go battery.publishPEMrequests()
	battery.logger.Info("Battery started\n")
	return battery
}

func (battery *Battery) setupBroker() broker.Broker {
	var brokerInstance, err = broker.NewBroker()
	if err != nil {
		battery.logger.Fatalf(err, "Broker instance could not be created")
	}
	return brokerInstance
}

func generateUuid() string {
	u := uuid.New()
	return u.String()
}

func (battery *Battery) handlePEMresponse(params ...[]byte) {
	battery.logger.Info("Received PEM response with the following params: %v", params)

	// todo verify index
	message := params[1]

	response := PEMResponse{}
	json.Unmarshal(message, &response)
	if response.BatteryID != battery.id {
		return
	}
	if response.ResponseType == GRANTED {
		battery.actOnGrantedRequest(response)
	} else if response.ResponseType == DENIED {
		battery.logger.Info("Request with id %s denied\n", response.ID)
	}
}

func (battery *Battery) publishPEMrequests() {
	for {
		time.Sleep(SENDING_INTERVAL_NS)
		for busy {
			time.Sleep(time.Second)
		}
		request := battery.getPEMRequest()
		battery.logger.Info("Sending %s request with id %s and battery id %s\n", request.RequestType, request.ID, request.BatteryID)

		jsonRequest, err := json.Marshal(request)
		if err != nil {
			battery.logger.Fatal(err)
		}
		battery.requests[request.ID] = request
		battery.broker.Publish(PEM_REQUESTS_TOPIC, request.BatteryID, string(jsonRequest))
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

func (battery *Battery) probabilisticallyCalculateRequest() PEMRequest {
	// based on state of charge send a request to charge or discharge the battery
	// TODO can be improved by using a more sophisticated algorithm
	if rand.Float64() < 0.5 {
		return battery.newRequest(CHARGE)
	}
	return battery.newRequest(DISCHARGE)
}

func (battery *Battery) actOnGrantedRequest(response PEMResponse) {
	battery.logger.Info("Request with id %s approved\n", response.ID)
	for busy {
		time.Sleep(1 * time.Second)
	}
	request := battery.requests[response.ID]

	if request.ID == "" {
		battery.logger.Info("Request with id %s not found\n", response.ID)
		return
	}
	if request.RequestType == CHARGE {
		battery.chargePacket()
	} else if request.RequestType == DISCHARGE {
		battery.dischargePacket()
	}
	battery.publishBatteryAction(request.RequestType)
	delete(battery.requests, request.ID)
}

func (battery *Battery) chargePacket() {
	battery.logger.Info("Charging packet of + %.2f kWh.\n", PACKET_KWH)
	battery.updateBattery(PACKET_KWH)
}

func (battery *Battery) dischargePacket() {
	battery.logger.Info("Discharging packet of - %.2f kWh.\n", PACKET_KWH)
	battery.updateBattery(-PACKET_KWH)
}

func (battery *Battery) updateBattery(chargeAmount float64) {
	busy = true
	currentBatteryCharge := battery.soc * BATTERY_CAPACITY_KWH
	currentBatteryCharge += chargeAmount
	battery.soc = (currentBatteryCharge / BATTERY_CAPACITY_KWH)
	time.Sleep(CHARGE_DISCHARGE_INTERVAL_MS * time.Millisecond)
	battery.logger.Info("After the update the new SoC is: %.4f\n", battery.soc)
	busy = false
}

func (battery *Battery) publishBatteryAction(actionType string) {
	action := BatteryAction{
		ID:         generateUuid(),
		BatteryID:  battery.id,
		ActionType: actionType,
	}
	json, err := json.Marshal(action)
	if err != nil {
		battery.logger.Fatal(err)
	}
	battery.broker.Publish(BATTERY_ACTIONS_TOPIC, battery.id, string(json))
}
