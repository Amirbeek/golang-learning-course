package main

import (
	"testing"
)

func TestProcessTruck(t *testing.T) {
	t.Run("processTruck", func(t *testing.T) {
		t.Run("should load and upload a truck cargo", func(t *testing.T) {
			nt := &NormalTruck{id: "Truck-1", cargo: 42}
			et := &ElectricTruck{id: "Electronic-Truck-1"}
			err := processTruck(nt)
			if err != nil {
				t.Fatalf("Err processing truck: %s", err)
			}
			err = processTruck(et)
			if err != nil {
				t.Fatalf("Err processing truck: %s", err)
			}
			//	Asserting
			//if nt.cargo != 0 {
			//	t.Fatalf("Normal Truck cargo should be 0: %d", nt.cargo)
			//}
			//if et.cargo != -1 {
			//	t.Fatalf("Electric Truck cargo should be 0: %d", et.cargo)
			//}
		})
	})
}
