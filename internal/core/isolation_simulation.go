package core

import (
	"context"
	"database/sql"
	"de/internal/storage/sqlstorage"
)

type SaleSimulation SimulationState[sqlstorage.Sale]

func SimulateLostUpdates(
	ctx context.Context,
	store *sqlstorage.Store,
) ([]SaleSimulation, error) {
	tx1, err := store.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
	})
	if err != nil {
		return nil, err
	}
	defer tx1.Rollback()

	tx2, err := store.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
	})
	if err != nil {
		return nil, err
	}
	defer tx2.Rollback()

	states := []SaleSimulation{{
		TxID:  "Explanation",
		Query: "A transaction re-executes a query returning a set of rows that satisfy a search condition and finds that the set of rows satisfying the condition has changed due to another recently-committed transaction.",
	}}
	tx1Sales, err := store.ListSales(ctx, tx1, 10, 0)
	if err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "1",
		Query: "SELECT * FROM sales",
		Rows:  tx1Sales,
	})

	if err := store.UpdateSaleQty(ctx, tx1, 1, 5); err != nil {
		return nil, err
	}

	tx1Sales, err = store.ListSales(ctx, tx1, 10, 0)
	if err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "1",
		Query: "UPDATE sales SET quantity = 5 WHERE id = 1",
		Rows:  tx1Sales,
	})

	tx2Sales, err := store.ListSales(ctx, tx2, 10, 0)
	if err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "2",
		Query: "SELECT * FROM sales",
		Rows:  tx2Sales,
	})

	if err := store.UpdateSaleQty(ctx, tx2, 1, 20); err != nil {
		return nil, err
	}
	tx2Sales, err = store.ListSales(ctx, tx2, 10, 0)
	if err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "2",
		Query: "UPDATE sales SET quantity = 20 WHERE id = 1",
		Rows:  tx2Sales,
	})

	if err := tx2.Commit(); err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "2",
		Query: "Commit Transaction",
	})

	tx1Sales, err = store.ListSales(ctx, tx1, 10, 0)
	if err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "1",
		Query: "SELECT * FROM sales",
		Rows:  tx1Sales,
	})
	states = append(states, SaleSimulation{
		TxID:  "1",
		Query: "Rollback Transaction",
	})

	return states, nil
}

func SimulatePhantomRead(
	ctx context.Context,
	store *sqlstorage.Store,
) ([]SaleSimulation, error) {
	tx1, err := store.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, err
	}
	defer tx1.Rollback()

	tx2, err := store.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, err
	}
	defer tx2.Rollback()

	states := []SaleSimulation{{
		TxID:  "Explanation",
		Query: "A transaction re-executes a query returning a set of rows that satisfy a search condition and finds that the set of rows satisfying the condition has changed due to another recently-committed transaction.",
	}}
	tx1Sales, err := store.ListSales(ctx, tx1, 10, 0)
	if err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "1",
		Query: "SELECT * FROM sales",
		Rows:  tx1Sales,
	})

	tx2Sales, err := store.ListSales(ctx, tx2, 10, 0)
	if err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "2",
		Query: "SELECT * FROM sales",
		Rows:  tx2Sales,
	})

	if err := store.InsertSale(ctx, tx2, 1, 10); err != nil {
		return nil, err
	}
	tx2Sales, err = store.ListSales(ctx, tx2, 10, 0)
	if err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "2",
		Query: "INSERT INTO sales(quantity, price) VALUES (10, 1)",
		Rows:  tx2Sales,
	})

	if err := tx2.Commit(); err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "2",
		Query: "Commit Transaction",
	})

	tx1Sales, err = store.ListSales(ctx, tx1, 10, 0)
	if err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "1",
		Query: "SELECT * FROM sales",
		Rows:  tx1Sales,
	})
	states = append(states, SaleSimulation{
		TxID:  "1",
		Query: "Rollback Transaction",
	})

	return states, nil
}

func SimulateNonRepeatableRead(
	ctx context.Context,
	store *sqlstorage.Store,
) ([]SaleSimulation, error) {
	tx1, err := store.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, err
	}
	defer tx1.Rollback()

	tx2, err := store.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, err
	}
	defer tx2.Rollback()

	states := []SaleSimulation{{
		TxID:  "Explanation",
		Query: "A transaction re-reads data it has previously read and finds that data has been modified by another transaction (that committed since the initial read).",
	}}
	tx1Sales, err := store.ListSales(ctx, tx1, 2, 0)
	if err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "1",
		Query: "SELECT * FROM sales",
		Rows:  tx1Sales,
	})

	tx2Sales, err := store.ListSales(ctx, tx2, 2, 0)
	if err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "2",
		Query: "SELECT * FROM sales",
		Rows:  tx2Sales,
	})

	if err := store.UpdateSaleQty(ctx, tx2, 1, 15); err != nil {
		return nil, err
	}
	tx2Sales, err = store.ListSales(ctx, tx2, 2, 0)
	if err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "2",
		Query: "UPDATE sales SET quantity = 15 WHERE id = 1",
		Rows:  tx2Sales,
	})

	if err := tx2.Commit(); err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "2",
		Query: "Commit Transaction",
	})

	tx1Sales, err = store.ListSales(ctx, tx1, 2, 0)
	if err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "1",
		Query: "SELECT * FROM sales",
		Rows:  tx1Sales,
	})
	states = append(states, SaleSimulation{
		TxID:  "1",
		Query: "Rollback Transaction",
	})

	return states, nil
}

func SimulateDirtyRead(
	ctx context.Context,
	store *sqlstorage.Store,
) ([]SaleSimulation, error) {
	tx1, err := store.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
	})
	if err != nil {
		return nil, err
	}
	defer tx1.Rollback()

	tx2, err := store.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
	})
	if err != nil {
		return nil, err
	}
	defer tx2.Rollback()

	states := []SaleSimulation{{
		TxID:  "Explanation",
		Query: "A transaction reads data written by a concurrent uncommitted transaction.",
	}}
	tx1Sales, err := store.ListSales(ctx, tx1, 2, 0)
	if err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "1",
		Query: "SELECT * FROM sales",
		Rows:  tx1Sales,
	})

	tx2Sales, err := store.ListSales(ctx, tx2, 2, 0)
	if err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "2",
		Query: "SELECT * FROM sales",
		Rows:  tx2Sales,
	})

	if err := store.UpdateSaleQty(ctx, tx2, 1, 15); err != nil {
		return nil, err
	}
	tx2Sales, err = store.ListSales(ctx, tx2, 2, 0)
	if err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "2",
		Query: "UPDATE sales SET quantity = 15 WHERE id = 1",
		Rows:  tx2Sales,
	})

	tx1Sales, err = store.ListSales(ctx, tx1, 2, 0)
	if err != nil {
		return nil, err
	}
	states = append(states, SaleSimulation{
		TxID:  "1",
		Query: "SELECT * FROM sales",
		Rows:  tx1Sales,
	})

	states = append(states, SaleSimulation{
		TxID:  "2",
		Query: "Rollback Transaction",
	})

	states = append(states, SaleSimulation{
		TxID:  "1",
		Query: "Rollback Transaction",
	})

	return states, nil
}
