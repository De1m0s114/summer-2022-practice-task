package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

type Trains []Train

type Train struct {
	TrainID            int
	DepartureStationID int
	ArrivalStationID   int
	Price              float32
	ArrivalTime        time.Time
	DepartureTime      time.Time
}
type Trains1 []Train1

type Train1 struct {
	TrainID            int
	DepartureStationID int
	ArrivalStationID   int
	Price              float32
	ArrivalTime        string
	DepartureTime      string
}

const longForm = "15:04:05"

func Parse() Trains {
	var trainsString Trains1

	byt, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(byt, &trainsString); err != nil {
		panic(err)
	}
	//time like time.Time
	trains := make([]Train, len(trainsString))

	for i, _ := range trainsString {
		trains[i].TrainID = trainsString[i].TrainID
		trains[i].DepartureStationID = trainsString[i].DepartureStationID
		trains[i].ArrivalStationID = trainsString[i].ArrivalStationID
		trains[i].Price = trainsString[i].Price
		trains[i].ArrivalTime, _ = time.Parse(longForm, trainsString[i].ArrivalTime)
		trains[i].DepartureTime, _ = time.Parse(longForm, trainsString[i].DepartureTime)
	}
	return trains
}

func main() {

	result, err := FindTrains("12", "1929", "price")

	if err != nil {
		fmt.Println(err)
	} else {
		for i, _ := range result {
			fmt.Printf("{TrainID: %v, DepartureStationID: %v, ArrivalStationID: %v, Price: %v, ArrivalTime: %v, DepartureTime: %v}\n", result[i].TrainID, result[i].DepartureStationID, result[i].ArrivalStationID, result[i].Price, result[i].ArrivalTime, result[i].DepartureTime)
		}
	}

}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {

	for i := 0; i < len(departureStation); i++ {
		if departureStation[i] < 48 || departureStation[i] > 57 {
			return nil, errors.New("bad departure station input")
		}
	}
	for i := 0; i < len(arrivalStation); i++ {
		if arrivalStation[i] < 48 || arrivalStation[i] > 57 {
			return nil, errors.New("bad arrival station input")
		}
	}
	intDepartureStation, err := strconv.Atoi(departureStation)
	if err != nil {
		return nil, errors.New("empty departure station")
	}
	intArrivalStation, err := strconv.Atoi(arrivalStation)
	if err != nil {
		return nil, errors.New("empty arrival station")
	}
	trains := Parse()
	var suitable_trains Trains

	for i, _ := range trains {
		if trains[i].DepartureStationID == intDepartureStation && trains[i].ArrivalStationID == intArrivalStation {
			suitable_trains = append(suitable_trains, trains[i])
		}
	}

	if criteria == "price" {
		for i := 0; i < len(suitable_trains); i++ {
			for j := i; j < len(suitable_trains); j++ {
				if suitable_trains[i].Price > suitable_trains[j].Price {
					suitable_trains[i], suitable_trains[j] = suitable_trains[j], suitable_trains[i]
				}
			}
		}
	} else if criteria == "arrival-time" {
		for i := 0; i < len(suitable_trains); i++ {
			for j := i; j < len(suitable_trains); j++ {
				if suitable_trains[j].ArrivalTime.Before(suitable_trains[i].ArrivalTime) {
					suitable_trains[i], suitable_trains[j] = suitable_trains[j], suitable_trains[i]
				}
			}
		}
	} else if criteria == "departure-time" {
		for i := 0; i < len(suitable_trains); i++ {
			for j := i; j < len(suitable_trains); j++ {
				if suitable_trains[j].DepartureTime.Before(suitable_trains[i].DepartureTime) {
					suitable_trains[i], suitable_trains[j] = suitable_trains[j], suitable_trains[i]
				}
			}
		}
	} else {
		return nil, errors.New("unsupported criteria")
	}
	if len(suitable_trains) > 3 {
		return suitable_trains[0:3], nil
	} else {
		return suitable_trains, nil
	}

}
