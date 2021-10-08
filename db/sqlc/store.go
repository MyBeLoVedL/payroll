package db

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type Store struct {
	q  *Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{q: New(db), db: db}
}

func (s *Store) execTrx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return errors.Errorf("tx error %v ,roll back error %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type UpdateBankParam struct {
	id                    int64
	bankName, bankAccount string
}

func (s *Store) UpdatePaymentMethodWithBank(ctx context.Context, arg UpdateBankParam) error {
	var err error
	err = s.execTrx(ctx, func(q *Queries) error {
		err = UpdatePayment("deposit", int64(arg.id))
		if err != nil {
			return err
		}

		err = q.InsertBank(context.Background(), InsertBankParams{arg.id, arg.bankName, arg.bankAccount})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
