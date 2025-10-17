package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type BankAccount struct {
	balance int
}

var (
	NotEnoughMoney = errors.New("not enough money")
)

func Deposit(b *BankAccount, amount int, wg *sync.WaitGroup) error {
	defer wg.Done()
	b.balance += amount
	fmt.Println("Depositing", b.balance)
	time.Sleep(1 * time.Second)
	b.balance += amount
	return nil
}
func Withdraw(b *BankAccount, amount int, wg *sync.WaitGroup) error {
	defer wg.Done()
	if b.balance < amount {
		return NotEnoughMoney
	}
	fmt.Println("Withdrawing", b.balance)
	time.Sleep(1 * time.Second)
	b.balance -= amount
	return nil
}

func main() {
	var wg sync.WaitGroup
	balance := BankAccount{1000}
	deposits := []int{200, 150, 300, 400}
	withdraws := []int{100, 25, 750, 2300, 1000, 300, 400}
	for _, deposit := range deposits {
		wg.Add(1)
		go func() {
			err := Deposit(&balance, deposit, &wg)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	for _, withdraw := range withdraws {
		wg.Add(1)
		go func() {
			err := Withdraw(&balance, withdraw, &wg)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	wg.Wait()
	fmt.Println("Balance:", balance.balance)
}
