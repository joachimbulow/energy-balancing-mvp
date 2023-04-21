package src

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"

	"github.com/joachimbulow/pem-energy-balance/src/broker"
)

const (
	PEM_REQUESTS_TOPIC          = "pem_requests"
	PEM_RESPONSES_TOPIC         = "pem_responses"
	FREQUENCY_MEASUREMENT_TOPIC = "frequency_measurements"
	BATTERY_ACTIONS_TOPIC       = "battery_actions"
	INERTIA_MEASUREMENT         = "inertia_measurements"
)

var (
	brokerInstance broker.Broker
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

	SENDING_INTERVAL_MS = 10000

	CHARGE_DISCHARGE_INTERVAL_MS = 10000
)

type Battery struct {
	ID             string
	BrokerInstance broker.Broker
	SoC            float64
	Requests       map[string]PEMRequest
}

var (
	battery = Battery{}
	busy    = false
)

func NewBattery() Battery {
	battery.ID = generateUuid()
	battery.BrokerInstance = setupBroker()
	go battery.BrokerInstance.Listen(PEM_RESPONSES_TOPIC, handlePEMresponse)
	go publishPEMrequests()
	return battery
}

func setupBroker() broker.Broker {
	var err error
	brokerInstance, err = broker.NewBroker(KAFKA)
	if err != nil {
		log.Print("Broker instance could not be created:", err)
	}
	return brokerInstance
}

func generateUuid() string {
	newUUID, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generated UUID:")
	fmt.Printf("%s", newUUID)
	return string(newUUID)
}

func handlePEMresponse(params ...[]byte) {
	fmt.Println("Received PEM response: ")
	for _, param := range params {
		fmt.Printf("'%s'", param)
	}

	// todo verify index
	message := params[1]

	response := PEMResponse{}
	json.Unmarshal(message, &response)
	if response.BatteryID != battery.ID {
		return
	}
	if response.ResponseType == GRANTED {
		actOnGrantedRequest(response)
	} else if response.ResponseType == DENIED {
		fmt.Printf("Request with id %s denied\n", response.ID)
	}
}

func publishPEMrequests() {
	for {
		time.Sleep(SENDING_INTERVAL_MS)
		for busy {
			time.Sleep(time.Second)
		}
		request := getPEMRequest()
		fmt.Printf("Sending %s request with id %s and battery id %s\n", request.RequestType, request.ID, request.BatteryID)

		jsonRequest, err := json.Marshal(request)
		if err != nil {
			log.Fatal(err)
		}
		battery.Requests[request.ID] = request
		battery.BrokerInstance.Publish(PEM_REQUESTS_TOPIC, request.BatteryID, string(jsonRequest))
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
	if battery.SoC == 0 {
		// first measurement should be normally distributed around soc mean and std dev
		battery.SoC = rand.NormFloat64()*SOC_STD + SOC_MEAN
	}
	fmt.Printf("Battery state of charge: %.2f\n", battery.SoC)
	return battery.SoC
}

func newRequest(requestType string) PEMRequest {
	return PEMRequest{
		ID:          generateUuid(),
		BatteryID:   battery.ID,
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
	fmt.Printf("Request with id %s approved\n", response.ID)
	request := battery.Requests[response.ID]

	if request.ID == "" {
		fmt.Printf("Request with id %s not found\n", response.ID)
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
	delete(battery.Requests, request.ID)
}

func chargePacket() {
	fmt.Printf("Charging packet of + %.2f kWh.\n", PACKET_KWH)
	updateBattery(PACKET_KWH)
}

func dischargePacket() {
	fmt.Printf("Discharging packet of - %.2f kWh.\n", PACKET_KWH)
	updateBattery(-PACKET_KWH)
}

func updateBattery(chargeAmount float64) {
	busy = true
	currentBatteryCharge := battery.SoC * BATTERY_CAPACITY_KWH
	currentBatteryCharge += chargeAmount
	battery.SoC = (currentBatteryCharge / BATTERY_CAPACITY_KWH)
	time.Sleep(CHARGE_DISCHARGE_INTERVAL_MS * time.Millisecond)
	fmt.Printf("After the update the new SoC is: %.4f\n", battery.SoC)
	busy = false
}

func publishBatteryAction(actionType string) {
	action := BatteryAction{
		ID:         generateUuid(),
		BatteryID:  battery.ID,
		ActionType: actionType,
	}
	json, err := json.Marshal(action)
	if err != nil {
		log.Fatal(err)
	}
	battery.BrokerInstance.Publish(BATTERY_ACTIONS_TOPIC, battery.ID, string(json))
}
