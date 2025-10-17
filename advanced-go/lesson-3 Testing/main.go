package main

import (
	"errors"
	"fmt"
)

type Truck interface {
	LoadCargo() error
	UnloadCargo() error
}

type NormalTruck struct {
	id    string
	cargo int
}

func (t *NormalTruck) Status() string {
	//TODO implement me
	panic("implement me")
}

type ElectricTruck struct {
	id      string
	cargo   int
	battery float64
}

func (e *ElectricTruck) Status() string {
	//TODO implement me
	panic("implement me")
}

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrTruckNotFound  = errors.New("truck not found")
)

func (t *NormalTruck) LoadCargo() error {
	t.cargo++
	return nil
}
func (t *NormalTruck) UnloadCargo() error {
	t.cargo--
	return nil
}

func (e *ElectricTruck) LoadCargo() error {
	e.cargo++
	e.battery--
	return nil
}
func (e *ElectricTruck) UnloadCargo() error {
	e.cargo += 0
	e.battery += -1
	return nil
}

func processTruck(truck Truck) error {
	fmt.Println("processTruck:", truck)

	err := truck.LoadCargo()
	if err != nil {
		return fmt.Errorf("error loading cargo: %v", err)
	}
	err = truck.UnloadCargo()
	if err != nil {
		return fmt.Errorf("error unloading cargo: %v", err)
	}
	return nil
}

func main() {

}
