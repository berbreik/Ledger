package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ledger/internal/domain"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) GetAll(ctx context.Context) ([]*domain.Account, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, owner_name, balance
		FROM accounts
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*domain.Account
	for rows.Next() {
		var account domain.Account
		if err := rows.Scan(&account.ID, &account.OwnerName, &account.Balance); err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *AccountRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM accounts
		WHERE id = $1
	`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf(domain.ErrAccountNotFound)
		}
		return err
	}
	return nil
}

func (r *AccountRepository) Create(ctx context.Context, account *domain.Account) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO accounts (id, owner_name, balance)
		VALUES ($1, $2, $3)
	`, account.ID, account.OwnerName, account.Balance)
	return err
}

func (r *AccountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, owner_name, balance
		FROM accounts
		WHERE id = $1
	`, id)

	var account domain.Account
	err := row.Scan(&account.ID, &account.OwnerName, &account.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &domain.Account{}, fmt.Errorf(domain.ErrAccountNotFound)
		}
		return &domain.Account{}, err
	}
	return &account, nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, id string, amount float64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
	`, amount, id)
	return err
}
