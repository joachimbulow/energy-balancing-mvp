package src

import (
	"encoding/json"
	"math/rand"
	"os/exec"
	"time"

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
}

var (
	battery = Battery{}
	busy    = false
)

func NewBattery() Battery {
	battery.id = generateUuid()
	battery.broker = setupBroker()
	battery.requests = make(map[string]PEMRequest)
	logger = util.NewLogger(battery.id)
	go battery.broker.Listen(PEM_RESPONSES_TOPIC, handlePEMresponse)
	go publishPEMrequests()
	logger.Info("Battery started\n")
	return battery
}

func setupBroker() broker.Broker {
	var brokerInstance, err = broker.NewBroker(KAFKA)
	if err != nil {
		logger.Fatalf(err, "Broker instance could not be created")
	}
	return brokerInstance
}

func generateUuid() string {
	newUUID, err := exec.Command("uuidgen").Output()
	if err != nil {
		logger.Fatal(err)
	}
	return string(newUUID)
}

func handlePEMresponse(params ...[]byte) {
	logger.Info("Received PEM response with the following params: %v", params)

	// todo verify index
	message := params[1]

	response := PEMResponse{}
	json.Unmarshal(message, &response)
	if response.BatteryID != battery.id {
		return
	}
	if response.ResponseType == GRANTED {
		actOnGrantedRequest(response)
	} else if response.ResponseType == DENIED {
		logger.Info("Request with id %s denied\n", response.ID)
	}
}

func publishPEMrequests() {
	for {
		time.Sleep(SENDING_INTERVAL_NS)
		for busy {
			time.Sleep(time.Second)
		}
		request := getPEMRequest()
		logger.Info("Sending %s request with id %s and battery id %s\n", request.RequestType, request.ID, request.BatteryID)

		jsonRequest, err := json.Marshal(request)
		if err != nil {
			logger.Fatal(err)
		}
		battery.requests[request.ID] = request
		battery.broker.Publish(PEM_REQUESTS_TOPIC, request.BatteryID, string(jsonRequest))
	}
}

func getPEMRequest() PEMRequest {
	stateOfCharge := measureSoC()

	var request PEMRequest
	if stateOfCharge < LowerBoundBatteryCapacity {
		request = newRequest(CHARGE)
	} else if stateOfCharge >= LowerBoundBatteryCapacity && stateOfCharge <= UpperBoundBatteryCapacity {
		request = probabilisticallyCalculateRequest()
	} else {
		request = newRequest(DISCHARGE)
	}
	return request
}

func measureSoC() float64 {
	if battery.soc == 0 {
		// first measurement should be normally distributed around soc mean and std dev
		battery.soc = rand.NormFloat64()*SOC_STD + SOC_MEAN
	}
	logger.Info("Battery state of charge: %.2f\n", battery.soc)
	return battery.soc
}

func newRequest(requestType string) PEMRequest {
	return PEMRequest{
		ID:          generateUuid(),
		BatteryID:   battery.id,
		RequestType: requestType,
	}
}

func probabilisticallyCalculateRequest() PEMRequest {
	// based on state of charge send a request to charge or discharge the battery
	// TODO can be improved by using a more sophisticated algorithm
	if rand.Float64() < 0.5 {
		return newRequest(CHARGE)
	}
	return newRequest(DISCHARGE)
}

func actOnGrantedRequest(response PEMResponse) {
	logger.Info("Request with id %s approved\n", response.ID)
	request := battery.requests[response.ID]

	if request.ID == "" {
		logger.Info("Request with id %s not found\n", response.ID)
		return
	}
	for busy {
		time.Sleep(1 * time.Second)
	}
	if request.RequestType == CHARGE {
		chargePacket()
	} else if request.RequestType == DISCHARGE {
		dischargePacket()
	}
	publishBatteryAction(request.RequestType)
	delete(battery.requests, request.ID)
}

func chargePacket() {
	logger.Info("Charging packet of + %.2f kWh.\n", PACKET_KWH)
	updateBattery(PACKET_KWH)
}

func dischargePacket() {
	logger.Info("Discharging packet of - %.2f kWh.\n", PACKET_KWH)
	updateBattery(-PACKET_KWH)
}

func updateBattery(chargeAmount float64) {
	busy = true
	currentBatteryCharge := battery.soc * BATTERY_CAPACITY_KWH
	currentBatteryCharge += chargeAmount
	battery.soc = (currentBatteryCharge / BATTERY_CAPACITY_KWH)
	time.Sleep(CHARGE_DISCHARGE_INTERVAL_MS * time.Millisecond)
	logger.Info("After the update the new SoC is: %.4f\n", battery.soc)
	busy = false
}

func publishBatteryAction(actionType string) {
	action := BatteryAction{
		ID:         generateUuid(),
		BatteryID:  battery.id,
		ActionType: actionType,
	}
	json, err := json.Marshal(action)
	if err != nil {
		logger.Fatal(err)
	}
	battery.broker.Publish(BATTERY_ACTIONS_TOPIC, battery.id, string(json))
}
