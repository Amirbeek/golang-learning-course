package main

import (
	"context"
	//"time"
)

func Process(ctx context.Context, accounts map[string]*Account) error {
	defer ctx.Done()

	return nil
}
func main() {
	//root := context.Background()
	//ctx, cancel := context.WithTimeout(root, 2*time.Second)
	//
	//acc1 := &Account{ID: "A1", Balance: 1000}
	//acc2 := &Account{ID: "A2", Balance: 500}
	//
	//t1 := DepositTransaction{"A1", 300}
	//t2 := TransferTransaction{"A1", "A2", 200}
	//t3 := WithdrawTransaction{"A2", 1000}

}
