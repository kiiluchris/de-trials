package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

type Store struct {
	DB *sql.DB
}

func NewStore(ctx context.Context) (*Store, error) {
	db, err := sql.Open("mysql", "detest:detest@tcp(localhost:3306)/detest")
	if err != nil {
		return nil, err
	}

	s := &Store{
		DB: db,
	}

	if err := s.init(ctx); err != nil {
		return nil, errors.Join(err, s.Close(ctx))
	}

	return s, nil
}

func (s *Store) Close(ctx context.Context) error {
	return s.DB.Close()
}

func (s *Store) init(ctx context.Context) error {
	s.DB.SetMaxOpenConns(25)
	s.DB.SetMaxIdleConns(25)
	s.DB.SetConnMaxIdleTime(5 * time.Minute)
	s.DB.SetConnMaxLifetime(5 * time.Minute)

	if err := s.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	if _, err := s.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS accounts (
			id INT AUTO_INCREMENT,
			balance INT NOT NULL,
			PRIMARY KEY (id)
		);
`); err != nil {
		return fmt.Errorf("create accounts table: %v", err)
	}

	if _, err := s.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS sales (
			id INT AUTO_INCREMENT,
			quantity INT NOT NULL,
			price INT NOT NULL,
			PRIMARY KEY (id)
		);
`); err != nil {
		return fmt.Errorf("create sales table: %v", err)
	}

	if _, err := s.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS employees (
			id INT AUTO_INCREMENT,
			name VARCHAR(300) NOT NULL,
			name2 VARCHAR(300) NOT NULL,
			PRIMARY KEY (id),
			INDEX(name)
		);
`); err != nil {
		return fmt.Errorf("create employees table: %v", err)
	}

	return s.Refresh(ctx)
}

func (s *Store) Refresh(ctx context.Context) error {
	if err := s.refreshAccounts(ctx); err != nil {
		return err
	}

	if err := s.refreshSales(ctx); err != nil {
		return err
	}

	if err := s.refreshEmployees(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Store) refreshAccounts(ctx context.Context) error {
	const truncateQuery = "TRUNCATE accounts"
	if _, err := s.DB.ExecContext(ctx, truncateQuery); err != nil {
		return fmt.Errorf("truncate accounts: %v", err)
	}

	const insertQuery = `
	INSERT INTO accounts(balance)
	VALUES (?);
	`
	var accid int
	for balance := 1000; balance > 0; balance -= 100 {
		accid += 1
		if _, err := s.DB.ExecContext(ctx, insertQuery, balance); err != nil {
			return fmt.Errorf("populate account %d of balance %d: %v", accid, balance, err)
		}
	}

	return nil
}

func (s *Store) refreshSales(ctx context.Context) error {
	const truncateQuery = "TRUNCATE sales"
	if _, err := s.DB.ExecContext(ctx, truncateQuery); err != nil {
		return fmt.Errorf("truncate accounts: %v", err)
	}

	if err := s.InsertSale(ctx, s.DB, 5, 10); err != nil {
		return fmt.Errorf("populate sale %d: %v", 1, err)
	}

	if err := s.InsertSale(ctx, s.DB, 4, 20); err != nil {
		return fmt.Errorf("populate sale %d: %v", 2, err)
	}

	return nil
}

func (s *Store) refreshEmployees(ctx context.Context) error {
	const truncateQuery = "TRUNCATE employees"
	if _, err := s.DB.ExecContext(ctx, truncateQuery); err != nil {
		return fmt.Errorf("truncate employees: %v", err)
	}

	const insertQuery = `
	INSERT INTO employees (name, name2)
	WITH RECURSIVE cte (n) AS (
		SELECT 1
		UNION ALL
		SELECT n + 1
		FROM cte WHERE n < 999
	)
	SELECT m.n, m.n FROM cte AS m, cte AS m1
	`
	if _, err := s.DB.ExecContext(ctx, insertQuery); err != nil {
		return fmt.Errorf("populate employees: %v", err)
	}

	count, err := s.CountEmployees(ctx)
	if err != nil {
		return fmt.Errorf("count employees: %w", err)
	}

	log.Printf("Inserted %d employee rows", count)

	return nil
}
