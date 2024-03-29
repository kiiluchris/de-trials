package sqlstorage

import (
	"context"
	"fmt"
)

type Sale struct {
	ID    uint64
	Price uint64
	Qty   uint64
}

func (s *Store) InsertSale(
	ctx context.Context,
	conn dbTx,
	price, qty uint64,
) error {
	const insertQuery = `
	INSERT INTO sales(quantity, price)
	VALUES (?, ?);
	`
	if _, err := conn.ExecContext(ctx, insertQuery, qty, price); err != nil {
		return fmt.Errorf("insert sale %d: %v", 1, err)
	}

	return nil
}

func (s *Store) UpdateSaleQty(
	ctx context.Context,
	conn dbTx,
	id, qty uint64,
) error {
	const query = "UPDATE sales SET quantity = ? WHERE id = ?"
	_, err := conn.ExecContext(ctx, query, qty, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) ListSales(
	ctx context.Context,
	conn dbTx,
	limit, offset uint64,
) ([]Sale, error) {
	const query = "SELECT id, quantity, price FROM sales LIMIT ? OFFSET ?"
	if limit == 0 {
		limit = 10
	}

	rows, err := conn.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accs []Sale
	for rows.Next() {
		var acc Sale
		if err := rows.Scan(&acc.ID, &acc.Qty, &acc.Price); err != nil {
			return nil, err
		}
		accs = append(accs, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accs, nil
}
