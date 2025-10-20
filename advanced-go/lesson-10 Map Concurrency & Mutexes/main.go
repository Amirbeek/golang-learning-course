package main

import (
	"errors"
	"sync"
)

var ErrTruckNotFound = errors.New("truck not found")

type FleetManager interface {
	AddTruck(id string, cargo int) error
	GetTruck(id string) (*Truck, error)
	RemoveTruck(id string) error
	UpdateTruckCargo(id string, cargo int) error
}

type Truck struct {
	ID    string
	Cargo int
}

type TruckManager struct {
	trucks map[string]*Truck
	sync.RWMutex
}

func NewTruckManager() TruckManager {
	return TruckManager{
		trucks: make(map[string]*Truck),
	}
}

func (tm *TruckManager) AddTruck(id string, cargo int) error {
	tm.Lock()
	defer tm.Unlock()
	_, exist := tm.trucks[id]
	if !exist {
		tm.trucks[id] = &Truck{id, cargo}
	}
	return ErrTruckNotFound
}

func (tm *TruckManager) GetTruck(id string) (*Truck, error) {
	tm.RLock()
	defer tm.RUnlock()

	truck, exist := tm.trucks[id]
	if exist {
		return truck, nil
	}
	return nil, ErrTruckNotFound
}

func (tm *TruckManager) RemoveTruck(id string) error {
	tm.Lock()
	defer tm.Unlock()
	_, exist := tm.trucks[id]
	if !exist {
		return ErrTruckNotFound
	}
	delete(tm.trucks, id)
	return nil
}

func (tm *TruckManager) UpdateTruckCargo(id string, cargo int) error {
	tm.Lock()
	defer tm.Unlock()
	truck, exist := tm.trucks[id]
	if !exist {
		return ErrTruckNotFound
	}
	truck.Cargo = cargo
	return nil
}
