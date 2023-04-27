package main

import (
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
	if nBatteriesEnv == "" {
		n := 1
		return n
	}
	nBatteries, err := strconv.Atoi(nBatteriesEnv)
	if err != nil {
		panic(err)
	}

	return nBatteries
}
