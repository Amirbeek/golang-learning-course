package main

import (
	"errors"
	"fmt"
)

type Truck interface {
	LoadCargo(c int) error
	UnloadCargo(c int) error
	Status() string
}

var (
	ErrNotEnoughCargo = errors.New("not enough cargo")
	ErrLowBattery     = errors.New("low battery")
	ErrLowFuel        = errors.New("low fuel")
)

type NormalTruck struct {
	id    string
	cargo int
}

type ElectricTruck struct {
	id      string
	cargo   int
	battery float64
}

type HybridTruck struct {
	id      string
	cargo   int
	fuel    int
	battery float64
}

// LoadCargo Normal Truck
func (t *NormalTruck) LoadCargo(c int) error {
	if c < 0 {
		return errors.New("negative load cargo")
	}
	t.cargo += c
	return nil
}
func (t *NormalTruck) UnloadCargo(c int) error {
	if c < 0 {
		return errors.New("negative unload cargo")
	}
	if c > t.cargo {
		return errors.New("cargo is too high")
	}
	t.cargo -= c
	return nil
}
func (t *NormalTruck) Status() string {
	return fmt.Sprintf("Normal truck:  %s | Cargo %d", t.id, t.cargo)
}

// LoadCargo Electric Truck
func (t *ElectricTruck) LoadCargo(c int) error {
	if c < 0 {
		return errors.New("negative unload cargo")
	}

	need := float64(t.cargo) * 2.0
	if t.battery < need {
		return ErrLowBattery
	}
	t.cargo += c
	t.battery += float64(t.cargo) * 100
	return nil
}
func (t *ElectricTruck) UnloadCargo(c int) error {
	if c < 0 {
		return errors.New("negative unload cargo")
	}
	if c > t.cargo {
		return errors.New("cargo is too high")
	}
	t.cargo -= c
	t.battery += -float64(t.cargo) * 100
	return nil
}
func (t *ElectricTruck) Status() string {
	return fmt.Sprintf("Electric Truck: %s | Cargo %d | Battery %.2f%%", t.id, t.cargo, t.battery)
}

// LoadCargo Hybrid Truck
func (h *HybridTruck) LoadCargo(c int) error {
	if c < 0 {
		return errors.New("negative unload cargo")
	}
	h.cargo += c
	h.battery += float64(c) * 100
	return nil
}
func (h *HybridTruck) UnloadCargo(c int) error {
	if c < 0 {
		return errors.New("negative unload cargo")
	}
	if c > h.cargo {
		return errors.New("cargo is too high")
	}
	need := float64(h.cargo) * 2.0
	if h.battery < need {
		return ErrLowBattery
	}
	h.cargo -= c
	h.battery -= float64(h.cargo) * 100
	h.fuel -= 10
	return nil
}
func (h *HybridTruck) Status() string {
	return fmt.Sprintf("Hybrid Truck: %s | Cargo %d | Fluid %dL | Battery %.2f%%", h.id, h.cargo, h.fuel, h.battery)
}

func processTrack(truck []Truck) error {
	for _, truck := range truck {
		fmt.Println("Processing truck:", truck)
		err := truck.LoadCargo(5)
		if err != nil {
			fmt.Println("Error loading cargo:", err)
		}
		err = truck.UnloadCargo(2)
		if err != nil {
			fmt.Println("Error unloading cargo:", err)
		}
		fmt.Println("After Processing truck:", truck.Status())
		fmt.Println("-----------")
	}
	return nil
}

func main() {
	normal := &NormalTruck{id: "N-1"}
	electric := &ElectricTruck{id: "E-1", battery: 100}
	hybrid := &HybridTruck{id: "H-1", fuel: 50, battery: 50}

	fleet := []Truck{normal, electric, hybrid}
	err := processTrack(fleet)
	if err != nil {
		fmt.Println("Error processing fleet:", err)
	}
}
