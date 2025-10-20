package main

import (
	"errors"
	"sync"
)

var (
	NotEnoughMoney error = errors.New("not enough money")
)

type Account struct {
	ID      string
	Balance float64
	mu      sync.Mutex
}

//type AccountManager struct {
//	accounts map[]
//}

func (ac *Account) DepositTransaction(amount float64) error {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.Balance += amount
	return nil
}
func (ac *Account) WithdrawTransaction(amount float64) error {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	if amount < ac.Balance {
		return NotEnoughMoney
	}
	ac.Balance -= amount
	return nil
}

func (ac *Account) TransferTransaction(receiver_id string, amount float64) error {

	return nil
}
