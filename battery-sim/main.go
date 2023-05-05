package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joachimbulow/pem-energy-balance/src"
)

var (
	nBatteries = getNumberOfBatteries()
)

func main() {
	initializeBatteries()
}

func initializeBatteries() {
	for i := 0; i < nBatteries; i++ {
		go startBattery()
	}

	for {
		time.Sleep(1 * time.Second)
	}
}

func startBattery() {
	src.NewBattery()
}

func getNumberOfBatteries() int {
	nBatteriesEnv := os.Getenv("N_BATTERIES")
	nBatteries, err := strconv.Atoi(nBatteriesEnv)
	if err != nil {
		// Print
		log.Println("Could not parse N_BATTERIES environment variable, defaulting to 2")
		nBatteries = 2
	}
	return nBatteries
}
