package main

import "context"

type Transaction interface {
	Process(ctx context.Context, accounts map[string]*Account) error
	Info() string
}
type DepositTransaction struct {
	Account string
	Amount  float64
}

type TransferTransaction struct {
	Sender_Account   string
	Reciever_Account string
	Amount           float64
}

type WithdrawTransaction struct {
	Account string
	Amount  float64
}
