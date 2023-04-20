package main

import (
	"fmt"

	"github.com/joachimbulow/pem-energy-balance/src"
	"github.com/joachimbulow/pem-energy-balance/src/broker"
)

const (
	PEMRequestsTopic     = "pem_requests"
	PEMResponsesTopic    = "pem_responses"
	FrequencyMeasurement = "frequency_measurements"
	BatteryActionsTopic  = "battery_actions"
	InertiaMeasurement   = "inertia_measurements"
)

var (
	// instance of the broker
	newBroker broker.Broker
)

type pemRequest struct {
	ID          string `json:"id"`
	RequestType string `json:"requestType"`
}

type pemResponse struct {
	ID           string `json:"id"`
	ResponseType string `json:"responseType"`
}

const (
	Charge    = "CHARGE"
	Discharge = "DISCHARGE"
)

const (
	Granted = "GRANTED"
	Denied  = "DENIED"
)

const (
	Redis = "REDIS"
	Kafka = "KAFKA"
)

var (
	nBatteries = 1
)

func main() {
	initialize()
}

func initialize() {
	for i := 0; i < nBatteries; i++ {
		go startBattery()
	}
}

func startBattery() {
	battery := src.NewBattery()
	fmt.Print("battery created: ", battery.ID)
}
