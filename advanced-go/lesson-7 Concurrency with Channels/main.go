package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Truck interface {
	LoadCargo() error
	UnloadCargo() error
}

type contextKey string

var UserIdKey contextKey = "userId"

type NormalTruck struct {
	id    string
	cargo int
}

type ElectricTruck struct {
	id      string
	cargo   int
	battery float64
}

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
	e.cargo--
	e.battery++
	return nil
}

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrTruckNotFound  = errors.New("truck not found")
)

func processTruck(ctx context.Context, truck Truck) error {
	fmt.Println("processTruck:", truck)
	userId := ctx.Value(UserIdKey).(int)
	log.Println("userId:", userId)

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	delay := time.Second * 1
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
		break
	}

	if err := truck.LoadCargo(); err != nil {
		return fmt.Errorf("error loading cargo: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	if err := truck.UnloadCargo(); err != nil {
		return fmt.Errorf("error unloading cargo: %v", err)
	}
	return nil
}

func processFleet(ctx context.Context, trucks []Truck) error {
	var wg sync.WaitGroup

	errorChan := make(chan error, len(trucks))

	for _, t := range trucks {
		wg.Add(1)
		go func(t Truck) {
			defer wg.Done()
			if err := processTruck(ctx, t); err != nil {
				fmt.Printf("error processing fleet: %v\n", err)
				errorChan <- err
			}

		}(t)
	}
	wg.Wait()
	close(errorChan)
	//select {
	//case err := <-errorChan:
	//	return err
	//default:
	//	return nil
	//}
	var errs []error

	for err := range errorChan {
		log.Printf("error processing fleet: %v\n", err)
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("some fleet processing failed: %v", errs)
	}

	return nil

}

func main() {

	ctx := context.Background()
	ctx = context.WithValue(ctx, UserIdKey, 442)

	fleet := []Truck{
		&NormalTruck{id: "NT1", cargo: 0},
		&ElectricTruck{id: "ET1", cargo: 0, battery: 100},
		&NormalTruck{id: "NT2", cargo: 0},
		&ElectricTruck{id: "ET2", cargo: 0, battery: 100},
	}
	start := time.Now()

	if err := processFleet(ctx, fleet); err != nil {
		fmt.Printf("Error processing fleet: %v", err)
	}
	fmt.Println("processTruck took:", time.Since(start))
	fmt.Println("All trucks processed successfully")
}
