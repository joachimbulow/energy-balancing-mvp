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
		// To spread out the spawn of the batteries so they spam less
		time.Sleep(2 * time.Second)
		go startBattery()
	}

	// To keep routines running' we start sleepin'
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
