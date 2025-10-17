package main

import (
	"errors"
	"fmt"
	"log"
)

type Truck struct {
	id string
}

func (t *Truck) Status() string {
	//TODO implement me
	panic("implement me")
}

type TruckError struct {
	Code int
	Msg  string
}

func (e *TruckError) Error() string {
	return fmt.Sprintf("Code %d: %s", e.Code, e.Msg)
}

func (t *Truck) LoadCargo() error {
	if t.id == "1" {
		return &TruckError{
			Code: 404,
			Msg:  "truck not found",
		}
	}
	fmt.Println("Loaded cargo:", t.id)
	return nil
}

func (t *Truck) UnloadCargo() error {
	if t.id == "2" {
		return &TruckError{
			Code: 500,
			Msg:  "unload failed",
		}
	}
	fmt.Println("Unloaded cargo:", t.id)
	return nil
}

func processTruck(t *Truck) error {
	fmt.Println("Processing truck:", t.id)

	if err := t.LoadCargo(); err != nil {
		return fmt.Errorf("load step failed: %w", err)
	}
	if err := t.UnloadCargo(); err != nil {
		return fmt.Errorf("unload step failed: %w", err)
	}

	fmt.Println("Completed:", t.id)
	return nil
}

func main() {
	trucks := []Truck{
		{id: "1"},
		{id: "2"},
		{id: "3"},
	}

	for i := range trucks {
		if err := processTruck(&trucks[i]); err != nil {
			var te *TruckError
			if errors.As(err, &te) {
				log.Printf("Truck %s error (code %d): %v\n", trucks[i].id, te.Code, err)
			} else {
				log.Printf("Truck %s error: %v\n", trucks[i].id, err)
			}
			continue
		}
	}
}
