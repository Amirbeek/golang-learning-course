package main

import (
	"errors"
	"fmt"
	"log"
)

type Truck struct {
	id string
}

var (
	ErrorNotImplemented = errors.New("not implemented")
	ErrorTrackNotFound  = errors.New("truck not found")
	ErrorUnloadCargo    = errors.New("unload cargo")
)

func (truck *Truck) LoadCargo() error {
	return ErrorTrackNotFound
}
func (truck *Truck) UnloadCargo() error {
	return ErrorUnloadCargo
}

func processTruck(truck Truck) error {
	fmt.Println("Truck is:", truck.id)
	if err := truck.LoadCargo(); err != nil {
		return fmt.Errorf("error loading cargo: %v", err)
	}
	if err := truck.UnloadCargo(); err != nil {

	}
	return ErrorNotImplemented
}

func main() {
	trucks := []Truck{
		{id: "Truck-1"},
		{id: "Truck-2"},
		{id: "Truck-3"},
	}

	for _, truck := range trucks {
		err := processTruck(truck)
		if err != nil {
			log.Printf("Error processing truck %s: %v\n", truck.id, err)
			continue // continue to the next truck instead of exiting
		}
	}
}
