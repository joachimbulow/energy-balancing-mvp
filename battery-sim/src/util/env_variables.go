package util

import (
	"log"
	"os"
	"strconv"
	"time"
)

func GetBrokerURL() string {
	if url := os.Getenv("BROKER_URL"); url != "" {
		return url
	}
	log.Println("BROKER_URL not set, using default 127.0.0.1:29092")
	return "127.0.0.1:29092" // Default to localhost
}

func GetBroker() string {
	if url := os.Getenv("BROKER"); url != "" {
		return url
	}
	log.Println("BROKER not set, using default KAFKA")
	return "KAFKA"
}

func GetNumberOfBatteries() int {
	if nBatteriesEnv := os.Getenv("N_BATTERIES"); nBatteriesEnv != "" {
		nBatteries, err := strconv.Atoi(nBatteriesEnv)
		if err == nil {
			return nBatteries
		}
	}
	log.Println("N_BATTERIES not set, using default 2")
	return 2 // Default to 2
}

func GetUpperBoundBatteryCapacity() float64 {
	if upperBound := os.Getenv("UPPER_BOUND_BATTERY_CAPACITY"); upperBound != "" {
		parsedValue, err := strconv.ParseFloat(upperBound, 64)
		if err == nil {
			return parsedValue
		}
	}
	log.Println("UPPER_BOUND_BATTERY_CAPACITY not set, using default 0.8")
	return 0.8 // Default to 0.8
}

func GetLowerBoundBatteryCapacity() float64 {
	if lowerBound := os.Getenv("LOWER_BOUND_BATTERY_CAPACITY"); lowerBound != "" {
		parsedValue, err := strconv.ParseFloat(lowerBound, 64)
		if err == nil {
			return parsedValue
		}
	}
	log.Println("LOWER_BOUND_BATTERY_CAPACITY not set, using default: 0.2")
	return 0.2 // Default to 0.2
}

func GetRequestInterval() time.Duration {
	if interval := os.Getenv("REQUEST_INTERVAL_SECONDS"); interval != "" {
		parsedValue, err := strconv.Atoi(interval)
		if err == nil {
			return time.Duration(parsedValue) * time.Second
		}
	}
	log.Println("REQUEST_INTERVAL_SECONDS not set, using default: 20 seconds")
	return 20 * time.Second // Default to 20 seconds
}

func GetPacketPowerW() int {
	if power := os.Getenv("PACKET_POWER_W"); power != "" {
		parsedValue, err := strconv.Atoi(power)
		if err == nil {
			return parsedValue
		}
	}
	log.Println("PACKET_POWER_W not set, using default: 4000 watts")
	return 4000 // Default to 4000 watts
}

func GetPacketTimeS() int {
	if time := os.Getenv("PACKET_TIME_S"); time != "" {
		parsedValue, err := strconv.Atoi(time)
		if err == nil {
			return parsedValue
		}
	}
	log.Println("PACKET_TIME_S not set, using default: 5 minutes")
	return 5 * 60 // Default to 5 minutes
}
