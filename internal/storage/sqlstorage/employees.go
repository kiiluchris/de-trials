package sqlstorage

import (
	"context"
	"fmt"
)

func (s *Store) CountEmployees(ctx context.Context) (uint64, error) {
	var count uint64
	row := s.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM employees")
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Store) AnalyzePaginationCAll(ctx context.Context) (string, error) {
	const query = `
	EXPLAIN ANALYZE SELECT *
	FROM employees
	WHERE id < 121452 ORDER BY id DESC
    LIMIT 10;
	`
	return s.analyzeQuery(ctx, query)
}

func (s *Store) AnalyzePaginationLOAll(ctx context.Context) (string, error) {
	const query = `
	EXPLAIN ANALYZE SELECT *
	FROM employees
	ORDER BY id DESC LIMIT 10 OFFSET 876550;
	`
	return s.analyzeQuery(ctx, query)
}

func (s *Store) AnalyzePaginationCNameName2(ctx context.Context) (string, error) {
	const query = `
	EXPLAIN ANALYZE SELECT id, name, name2
	FROM employees
	WHERE id < 121452 ORDER BY id DESC
    LIMIT 10;
	`
	return s.analyzeQuery(ctx, query)
}

func (s *Store) AnalyzePaginationLONameName2(ctx context.Context) (string, error) {
	const query = `
	EXPLAIN ANALYZE SELECT id, name, name2
	FROM employees
	ORDER BY id DESC LIMIT 10 OFFSET 876550;
	`
	return s.analyzeQuery(ctx, query)
}

func (s *Store) AnalyzePaginationCName(ctx context.Context) (string, error) {
	const query = `
	EXPLAIN ANALYZE SELECT id, name
	FROM employees
	WHERE id < 121452 ORDER BY id DESC
    LIMIT 10;
	`
	return s.analyzeQuery(ctx, query)
}

func (s *Store) AnalyzePaginationLOName(ctx context.Context) (string, error) {
	const query = `
	EXPLAIN ANALYZE SELECT id, name
	FROM employees
	ORDER BY id DESC LIMIT 10 OFFSET 876550;
	`
	return s.analyzeQuery(ctx, query)
}

func (s *Store) analyzeQuery(
	ctx context.Context,
	query string,
) (string, error) {
	row := s.DB.QueryRowContext(ctx, query)
	var result string
	if err := row.Scan(&result); err != nil {
		return "", fmt.Errorf("explain analyze * employee: %v", err)
	}

	return result, nil
}

func (s *Store) AnalyzePaginationCursor(
	ctx context.Context,
) (string, error) {
	const query = `
	EXPLAIN ANALYZE SELECT id FROM employees WHERE id < 121452 ORDER BY id DESC
    LIMIT 10;
	`
	row := s.DB.QueryRowContext(ctx, query)
	var result string
	if err := row.Scan(&result); err != nil {
		return "", fmt.Errorf("explain analyze * employee: %v", err)
	}

	return result, nil
}

func (s *Store) AnalyzePaginationLimitOffset(
	ctx context.Context,
) (string, error) {
	const query = `
	EXPLAIN ANALYZE SELECT id FROM employees ORDER BY id DESC LIMIT 10 OFFSET 876550;
	`
	row := s.DB.QueryRowContext(ctx, query)
	var result string
	if err := row.Scan(&result); err != nil {
		return "", fmt.Errorf("explain analyze * employee: %v", err)
	}

	return result, nil
}

func (s *Store) AnalyzeAllColumnSelect(
	ctx context.Context,
) (string, error) {
	const query = `
	EXPLAIN ANALYZE SELECT * FROM employees WHERE id = 777;
	`
	row := s.DB.QueryRowContext(ctx, query)
	var result string
	if err := row.Scan(&result); err != nil {
		return "", fmt.Errorf("explain analyze * employee: %v", err)
	}

	return result, nil
}

func (s *Store) AnalyzePrimaryKeyPlusIndexSelect(
	ctx context.Context,
) (string, error) {
	const query = `
	EXPLAIN ANALYZE SELECT id, name FROM employees where id = 777;
	`
	row := s.DB.QueryRowContext(ctx, query)
	var result string
	if err := row.Scan(&result); err != nil {
		return "", fmt.Errorf("explain analyze employee: %v", err)
	}

	return result, nil
}

func (s *Store) AnalyzeAllExplicitIndexSelect(
	ctx context.Context,
) (string, error) {
	const query = `
	EXPLAIN ANALYZE SELECT id, name, name2 FROM employees where id = 777;
	`
	row := s.DB.QueryRowContext(ctx, query)
	var result string
	if err := row.Scan(&result); err != nil {
		return "", fmt.Errorf("explain analyze employee: %v", err)
	}

	return result, nil
}

func (s *Store) AnalyzePrimaryKeySelect(
	ctx context.Context,
) (string, error) {
	const query = `
	EXPLAIN ANALYZE SELECT id FROM employees where id = 777;
	`
	row := s.DB.QueryRowContext(ctx, query)
	var result string
	if err := row.Scan(&result); err != nil {
		return "", fmt.Errorf("explain analyze pk employee: %v", err)
	}

	return result, nil
}

func (s *Store) AnalyzeIndexedColumnSelect(
	ctx context.Context,
) (string, error) {
	const query = `
	EXPLAIN ANALYZE SELECT name FROM employees WHERE name = 777;
	`
	row := s.DB.QueryRowContext(ctx, query)
	var result string
	if err := row.Scan(&result); err != nil {
		return "", fmt.Errorf("explain analyze pk employee: %v", err)
	}

	return result, nil
}

func (s *Store) AnalyzeUnindexedColumnSelect(
	ctx context.Context,
) (string, error) {
	const query = `
	EXPLAIN ANALYZE SELECT name2 FROM employees where name2 = 777;
	`
	row := s.DB.QueryRowContext(ctx, query)
	var result string
	if err := row.Scan(&result); err != nil {
		return "", fmt.Errorf("explain analyze unindexed employee: %v", err)
	}

	return result, nil
}
