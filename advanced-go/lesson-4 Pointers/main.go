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
	return fmt.Sprintf("NormalTruck %s | cargo=%d", t.id, t.cargo)
}

type ElectricTruck struct {
	id      string
	cargo   int
	battery float64
}

func (e *ElectricTruck) Status() string {
	return fmt.Sprintf("ElectricTruck %s | cargo=%d | battery=%.2f", e.id, e.cargo, e.battery)
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
	if t.cargo > 0 {
		t.cargo--
	}
	return nil
}

func (e *ElectricTruck) LoadCargo() error {
	e.cargo++
	e.battery--
	return nil
}
func (e *ElectricTruck) UnloadCargo() error {
	if e.cargo > 0 {
		e.cargo--
	}
	e.battery--
	return nil
}

func UnloadCargo(e *ElectricTruck) error {
	e.cargo = 0
	e.battery += -1
	return nil
}

func processTruck(truck Truck) error {
	fmt.Println("processTruck:", truck)

	if err := truck.LoadCargo(); err != nil {
		return fmt.Errorf("error loading cargo: %w", err)
	}
	if err := truck.UnloadCargo(); err != nil {
		return fmt.Errorf("error unloading cargo: %w", err)
	}
	return nil
}

func point(foo *string) {
	fmt.Println("point:", *foo)
}

func main() {
	foo := "amir"
	point(&foo)
	fmt.Println(&foo)

}
