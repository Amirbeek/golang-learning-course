package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
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

func processFleet(truck Truck) error {
	fmt.Println("processTruck:", truck)

	err := truck.LoadCargo()
	if err != nil {
		return fmt.Errorf("error loading cargo: %v", err)
	}

	time.Sleep(500 * time.Millisecond) // 🕒 0.5 sekund kutish
	err = truck.UnloadCargo()
	if err != nil {
		return fmt.Errorf("error unloading cargo: %v", err)
	}
	return nil
}
func processTruck(trucks []Truck) error {
	var wg sync.WaitGroup

	for _, t := range trucks {
		wg.Add(1)
		go func(t Truck) {
			defer wg.Done()
			err := processFleet(t)
			if err != nil {
				fmt.Printf("error processing fleet: %v\n", err)
			}
		}(t)
	}

	wg.Wait()
	return nil
}

func main() {
	fleet := []Truck{
		&NormalTruck{id: "NT1", cargo: 0},
		&ElectricTruck{id: "ET1", cargo: 0, battery: 100},
		&NormalTruck{id: "NT2", cargo: 0},
		&ElectricTruck{id: "ET2", cargo: 0, battery: 100},
	}
	start := time.Now()

	if err := processTruck(fleet); err != nil {
		fmt.Printf("Error processing fleet: %v", err)
	}
	fmt.Println("processTruck took:", time.Since(start))

	fmt.Println("All Truck processed successfully")
}
