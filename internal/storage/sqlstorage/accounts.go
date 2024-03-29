package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
)

type dbTx interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

func (s *Store) FailedAtomicTransfer(
	ctx context.Context,
	from, to uint64,
	amount uint64,
) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.failedTransfer(ctx, tx, from, to, amount); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) FailedNonAtomicTransfer(
	ctx context.Context,
	from, to uint64,
	amount uint64,
) error {
	return s.failedTransfer(ctx, s.DB, from, to, amount)
}

func (s *Store) NonAtomicTransfer(
	ctx context.Context,
	from, to uint64,
	amount uint64,
) error {
	return s.transfer(ctx, s.DB, from, to, amount)
}

func (s *Store) AtomicTransfer(
	ctx context.Context,
	from, to uint64,
	amount uint64,
) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.transfer(ctx, tx, from, to, amount); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) failedTransfer(
	ctx context.Context,
	conn dbTx,
	from, to uint64,
	amount uint64,
) error {
	balance, err := s.getBalance(ctx, conn, from)
	if err != nil {
		return err
	}

	if balance < amount {
		return fmt.Errorf("account balance is less than transfer amount %d", amount)
	}

	if err := s.withdrawAmount(ctx, conn, from, amount); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	cancel()

	if err := s.depositAmount(ctx, conn, to, amount); err != nil {
		return err
	}

	return nil
}

func (s *Store) transfer(
	ctx context.Context,
	conn dbTx,
	from, to uint64,
	amount uint64,
) error {
	balance, err := s.getBalance(ctx, conn, from)
	if err != nil {
		return err
	}

	if balance < amount {
		return fmt.Errorf("account balance is less than transfer amount %d", amount)
	}

	if err := s.withdrawAmount(ctx, conn, from, amount); err != nil {
		return err
	}

	if err := s.depositAmount(ctx, conn, to, amount); err != nil {
		return err
	}

	return nil
}

func (s *Store) depositAmount(
	ctx context.Context,
	conn dbTx,
	id uint64,
	amount uint64,
) error {
	const query = "UPDATE accounts SET balance = balance + ? WHERE id = ?"
	_, err := conn.ExecContext(ctx, query, amount, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) withdrawAmount(
	ctx context.Context,
	conn dbTx,
	id uint64,
	amount uint64,
) error {
	const query = "UPDATE accounts SET balance = balance - ? WHERE id = ?"
	_, err := conn.ExecContext(ctx, query, amount, id)
	if err != nil {
		return err
	}

	return nil
}

type Account struct {
	ID      uint64
	Balance uint64
}

func (s *Store) ListAccounts(
	ctx context.Context,
	limit, offset uint64,
) ([]Account, error) {
	const query = "SELECT id, balance FROM accounts LIMIT ? OFFSET ?"
	if limit == 0 {
		limit = 10
	}

	rows, err := s.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accs []Account
	for rows.Next() {
		var acc Account
		if err := rows.Scan(&acc.ID, &acc.Balance); err != nil {
			return nil, err
		}
		accs = append(accs, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accs, nil
}

func (s *Store) getBalance(
	ctx context.Context,
	conn dbTx,
	id uint64,
) (uint64, error) {
	const query = "SELECT balance FROM accounts WHERE id = ? LIMIT 1"
	row := conn.QueryRowContext(ctx, query, id)
	var balance uint64
	if err := row.Scan(&balance); err != nil {
		return 0, err
	}

	return balance, nil
}
