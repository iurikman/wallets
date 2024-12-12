package models

import "errors"

var (
	ErrBalanceBelowZero        = errors.New("balance is below zero")
	ErrWalletNotFound          = errors.New("wallet not found")
	ErrWalletIDIsEmpty         = errors.New("wallet ID is empty")
	ErrAmountIsZero            = errors.New("amount is zero")
	ErrTransactionTypeIsEmpty  = errors.New("transaction type is empty")
	ErrChangeBalanceData       = errors.New("change balance data is wrong")
	ErrOperationTypeNotAllowed = errors.New("operation type not allowed")
)
