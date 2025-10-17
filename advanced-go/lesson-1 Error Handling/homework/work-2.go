package main

//
//import (
//	"errors"
//	"fmt"
//	"log"
//)
//
//type Truck struct {
//	id string
//}
//
//// Sentinel errors (Go style: ErrXxx)
//var (
//	ErrCargoMissing      = errors.New("cargo missing")
//	ErrDestinationClosed = errors.New("destination closed")
//)
//
//func (t *Truck) LoadCargo() error {
//	if t.id == "1" {
//		return ErrCargoMissing
//	}
//	fmt.Println("Loaded cargo:", t.id)
//	return nil
//}
//
//func (t *Truck) DeliverCargo() error {
//	if t.id == "2" {
//		return ErrDestinationClosed
//	}
//	fmt.Println("Delivered cargo:", t.id)
//	return nil
//}
//
//func processDelivery(t *Truck) error {
//	fmt.Println("Processing truck:", t.id)
//
//	if err := t.LoadCargo(); err != nil {
//		return fmt.Errorf("load step failed: %w", err)
//	}
//	if err := t.DeliverCargo(); err != nil {
//		return fmt.Errorf("deliver step failed: %w", err)
//	}
//
//	fmt.Println("Completed:", t.id)
//	return nil
//}
//
//func main() {
//	trucks := []Truck{
//		{id: "1"},
//		{id: "2"},
//		{id: "3"},
//	}
//
//	for i := range trucks {
//		if err := processDelivery(&trucks[i]); err != nil {
//			switch {
//			case errors.Is(err, ErrCargoMissing):
//				log.Printf("Truck %s: cannot start, %v\n", trucks[i].id, err)
//			case errors.Is(err, ErrDestinationClosed):
//				log.Printf("Truck %s: delivery blocked, %v\n", trucks[i].id, err)
//			default:
//				log.Printf("Truck %s: %v\n", trucks[i].id, err)
//			}
//			continue
//		}
//	}
//}
