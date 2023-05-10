package db

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type db struct {
	pool *pgxpool.Pool
}

// Instance - global database instance
var Instance = &db{}

// InUse - Use this getter to figure out if events should be written to database
func (db *db) InUse() bool {
	return db.pool != nil
}

func (db *db) QueryRows(query string, args ...interface{}) (rows pgx.Rows, err error) {
	rows, err = db.pool.Query(context.Background(), query, args...)
	return
}

func (db *db) QueryRow(query string, args ...interface{}) pgx.Row {
	//    defer SentryRecover()
	return db.pool.QueryRow(context.Background(), query, args...)
}

func (db *db) Execute(query string, args ...interface{}) (tag pgconn.CommandTag, err error) {
	//    defer SentryRecover()
	tag, err = db.pool.Exec(context.Background(), query, args...)
	return
}
