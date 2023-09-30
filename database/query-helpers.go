package database

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

// QueryResults is a wrapper for a database query result
type QueryResults[T interface{}] struct {
	Err   error
	Value []T
}

type ScalarResult[T interface{}] struct {
	Err   error
	Value T
}

// NonQueryResult is a wrapper for a database query report
type NonQueryResult struct {
	Err    error
	Result sql.Result
}

// ExecuteQuery is a wrapper for a database query
func ExecuteQuery[T any](db *sqlx.DB, query string, args ...interface{}) chan QueryResults[T] {
	resultChan := make(chan QueryResults[T], 1)

	go func() {
		var results []T
		err := db.Select(&results, query, args...)

		resultChan <- QueryResults[T]{
			Err:   err,
			Value: results,
		}
		close(resultChan)
	}()

	return resultChan
}

func ExecuteNamedQuery[T any](db *sqlx.DB, query string, args interface{}) chan QueryResults[T] {
	resultChan := make(chan QueryResults[T], 1)
	go func() {
		var result QueryResults[T]
		var rows []T
		err := db.Select(&rows, db.Rebind(query), args)
		if err != nil {
			result.Err = err
		} else {
			result.Value = rows
		}
		resultChan <- result
		close(resultChan)
	}()
	return resultChan
}

// ExecuteNonQuery is a wrapper for a database query that does not return a value
func ExecuteNonQuery(db *sqlx.DB, query string, args ...interface{}) chan NonQueryResult {
	resultChan := make(chan NonQueryResult, 1)

	go func() {
		result, err := db.Exec(query, args...)
		queryResult := NonQueryResult{
			Err:    err,
			Result: result,
		}
		resultChan <- queryResult
		close(resultChan)
	}()

	return resultChan
}

// ExecuteNamedNonQuery is a wrapper for a database query
func ExecuteNamedNonQuery(db *sqlx.DB, query string, args interface{}) chan NonQueryResult {
	resultChan := make(chan NonQueryResult, 1)
	go func() {
		res, err := db.NamedExec(query, args)
		resultChan <- NonQueryResult{
			Err:    err,
			Result: res,
		}
		close(resultChan)
	}()
	return resultChan
}

// ExecuteScalar is a wrapper for a database query that returns a single value
func ExecuteScalar[T any](db *sqlx.DB, query string, args ...interface{}) chan ScalarResult[T] {
	resultChan := make(chan ScalarResult[T], 1)

	go func() {
		var result T
		err := db.Get(&result, query, args...)

		resultChan <- ScalarResult[T]{
			Err:   err,
			Value: result,
		}
		close(resultChan)
	}()

	return resultChan
}

// ExecuteNamedScalar is a wrapper for a database query that returns a single value
func ExecuteNamedScalar[T any](db *sqlx.DB, query string, data interface{}) chan ScalarResult[T] {
	resultChan := make(chan ScalarResult[T], 1)

	go func() {
		var result T
		err := db.Get(&result, query, data)

		resultChan <- ScalarResult[T]{
			Err:   err,
			Value: result,
		}
		close(resultChan)
	}()

	return resultChan
}
