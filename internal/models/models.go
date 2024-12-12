package models

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID        uuid.UUID
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   bool
}

type Transaction struct {
	TransactionID uuid.UUID `json:"id"`
	WalletID      uuid.UUID `json:"walletId"`
	Amount        float64   `json:"amount"`
	OperationType string    `json:"transactionType"`
	ExecutedAt    time.Time `json:"executedAt"`
}

func (t Transaction) Validate() error {
	if t.WalletID == uuid.Nil {
		return ErrWalletIDIsEmpty
	}

	if _, ok := allowedOperationTypes[t.OperationType]; !ok {
		return ErrOperationTypeNotAllowed
	}

	if t.Amount <= 0 {
		return ErrAmountIsZero
	}

	if t.OperationType == "" {
		return ErrTransactionTypeIsEmpty
	}

	return nil
}

//nolint:gochecknoglobals
var allowedOperationTypes = map[string]struct{}{
	"DEPOSIT":  {},
	"WITHDRAW": {},
}
