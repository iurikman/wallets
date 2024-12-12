package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/iurikman/wallets/internal/models"
)

type db interface {
	CreateWallet(ctx context.Context) (*models.Wallet, error)
	GetWallet(ctx context.Context, id uuid.UUID) (*models.Wallet, error)
	Deposit(ctx context.Context, transaction models.Transaction) error
	Withdraw(ctx context.Context, transaction models.Transaction) error
}

type Service struct {
	db db
}

func New(db db) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) CreateWallet(ctx context.Context) (*models.Wallet, error) {
	createdWallet, err := s.db.CreateWallet(ctx)
	if err != nil {
		return nil, fmt.Errorf("s.db.CreateWallet(ctx) err: %w", err)
	}

	return createdWallet, nil
}

func (s *Service) GetWallet(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	wallet, err := s.db.GetWallet(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("s.db.GetWallet(ctx, id) err: %w", err)
	}

	return wallet, nil
}

func (s *Service) Withdraw(ctx context.Context, transaction models.Transaction) error {
	err := s.db.Withdraw(ctx, transaction)
	if err != nil {
		return fmt.Errorf("s.db.Withdraw() err: %w", err)
	}

	return nil
}

func (s *Service) Deposit(ctx context.Context, transaction models.Transaction) error {
	err := s.db.Deposit(ctx, transaction)
	if err != nil {
		return fmt.Errorf("s.db.Deposit() err: %w", err)
	}

	return nil
}
