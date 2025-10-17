package main

import (
	"errors"
	"fmt"
	"log"
)

type Car struct {
	id   string
	fuel int
}

func (car Car) StartEngine() error {
	if car.fuel == 0 {
		return errors.New("fuel is empty")
	} else {
		fmt.Println("engine started")
	}
	return nil
}

func processStartCar(car Car) error {
	fmt.Println("Starting ", car.id)
	return nil
}
func main() {
	cars := []Car{
		{id: "1", fuel: 0},
		{id: "2", fuel: 2},
		{id: "3", fuel: 1},
	}
	for _, car := range cars {
		_ = processStartCar(car)
		if err := car.StartEngine(); err != nil {
			log.Printf("Error: %v", err)
			continue
		}
	}

}
