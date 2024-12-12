package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iurikman/wallets/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	log "github.com/sirupsen/logrus"
)

func (p *Postgres) CreateWallet(ctx context.Context) (*models.Wallet, error) {
	createdWallet := new(models.Wallet)

	timeNow := time.Now()

	query := `INSERT INTO wallets (id, balance, created_at, updated_at, deleted) 
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id, balance, created_at, updated_at, deleted
				`

	if err := p.db.QueryRow(
		ctx,
		query,
		uuid.New(),
		0,
		timeNow,
		timeNow,
		false,
	).Scan(
		&createdWallet.ID,
		&createdWallet.Balance,
		&createdWallet.CreatedAt,
		&createdWallet.UpdatedAt,
		&createdWallet.Deleted,
	); err != nil {
		return nil, fmt.Errorf("creating wallet error: %w", err)
	}

	return createdWallet, nil
}

func (p *Postgres) GetWallet(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet

	query := `	SELECT id, balance, created_at, updated_at, deleted 
				FROM wallets 
				WHERE id = $1 AND deleted = false`

	err := p.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&wallet.ID,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
		&wallet.Deleted,
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, models.ErrWalletNotFound
	case err != nil:
		return nil, fmt.Errorf("getting wallet by id error: %w", err)
	}

	return &wallet, nil
}

func (p *Postgres) Deposit(ctx context.Context, transaction models.Transaction) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("p.db.Begin(ctx) err: %w", err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Warnf("deposit tx.Rollback(ctx) err: %v", err)
		}
	}()

	err = p.updateWalletBalance(ctx, tx, transaction.WalletID, transaction.Amount)
	if err != nil {
		return models.ErrChangeBalanceData
	}

	err = p.saveTransaction(ctx, tx, transaction)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction commit err: %w", err)
	}

	return nil
}

func (p *Postgres) Withdraw(ctx context.Context, transaction models.Transaction) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("p.db.Begin(ctx) err: %w", err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Warnf("withdraw tx.Rollback(ctx) err: %v", err)
		}
	}()

	err = p.updateWalletBalance(ctx, tx, transaction.WalletID, -transaction.Amount)

	switch {
	case errors.Is(err, models.ErrBalanceBelowZero):
		return models.ErrBalanceBelowZero
	case err != nil:
		return models.ErrChangeBalanceData
	}

	err = p.saveTransaction(ctx, tx, transaction)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction commit err: %w", err)
	}

	return nil
}

func (p *Postgres) updateWalletBalance(ctx context.Context, tx pgx.Tx, walletID uuid.UUID, amount float64) error {
	query := `	UPDATE wallets SET balance = balance + $2, updated_at = $3
                WHERE id = $1 and deleted = false 
				RETURNING id, balance
				`

	_, err := tx.Exec(
		ctx,
		query,
		walletID,
		amount,
		time.Now(),
	)

	var pgErr *pgconn.PgError

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.ErrWalletNotFound
	case errors.As(err, &pgErr) && pgErr.Code == pgerrcode.CheckViolation:
		return models.ErrBalanceBelowZero
	case err != nil:
		return fmt.Errorf("updating wallet error: %w", err)
	}

	return nil
}

func (p *Postgres) saveTransaction(ctx context.Context, tx pgx.Tx, transaction models.Transaction) error {
	var executedOperation models.Transaction

	query := `INSERT INTO transactions_history
    (id, wallet_id, amount, transaction_type, executed_at)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, wallet_id, amount, transaction_type, executed_at`

	err := tx.QueryRow(
		ctx,
		query,
		uuid.New(),
		transaction.WalletID,
		transaction.Amount,
		transaction.OperationType,
		time.Now(),
	).Scan(
		&executedOperation.TransactionID,
		&executedOperation.WalletID,
		&executedOperation.Amount,
		&executedOperation.OperationType,
		&executedOperation.ExecutedAt,
	)

	var pgErr *pgconn.PgError

	switch {
	case errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation:
		return models.ErrWalletNotFound
	case err != nil:
		return fmt.Errorf("transaction writing to database err: %w", err)
	}

	return nil
}
