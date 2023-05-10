package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/k0rean-rand0m/event-horizon/src/config"
	"log"
)

// Solidity to Postgres types mapping
var typesMapping = map[string]string{
	"address": "text",
	"bool":    "boolean",
	"int8":    "smallint",
	"int16":   "smallint",
	"int32":   "integer",
	"int64":   "bigint",
	"int256":  "text",
	"uint8":   "smallint",
	"uint16":  "smallint",
	"uint32":  "integer",
	"uint64":  "text",
	"uint256": "text",
	"string":  "text",
	"bytes":   "bytea",
}

func (db *db) Setup() error {

	// Connecting to a database
	log.Println("DB: connecting to database")
	connection, err := pgxpool.Connect(context.Background(), config.Config.DbUrl)
	if err != nil {
		return err
	}
	db.pool = connection
	log.Println("DB: connected to database")

	// Introspecting database

	log.Println("DB: introspecting database")
	// Checks that "event_horizon" schema exists
	var exists bool
	err = db.QueryRow(`
		select exists(select 1 from information_schema.schemata where schema_name = 'event_horizon');
	`).Scan(&exists)
	if err != nil {
		return err
	}

	// Creating "event_horizon" schema if it doesn't exist
	if !exists {
		log.Println("DB: creating \"event_horizon\" schema")
		_, err = db.Execute(`
			create schema event_horizon;
		`)
		if err != nil {
			return err
		}
	}

	// Checks that table exists for each event. If not, creates it.
	for _, event := range config.Config.Events {

		if event.Table == "" {
			return errors.New("table name for " + event.Label + "is not specified")
		}

		exists = false
		err = db.QueryRow(
			`
			select EXISTS (
			   select from information_schema.tables 
			   where  table_schema = 'event_horizon'
			   and    table_name   = $1
		   );`, event.Table,
		).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			log.Println("DB: creating table", event.Table)
			err = createEventTable(event)
			if err != nil {
				return err
			}
		} else {
			// ToDo: check that table structure is correct
		}
	}

	return nil
}

func createEventTable(event config.Event) error {
	query := "create table event_horizon." + event.Table + " (hash text, network text, "
	for _, argument := range event.Arguments {
		if !argument.Indexed {
			continue
		}
		query += "\"" + argument.Label + "\" " + typesMapping[argument.Type] + " not null, "
	}
	query += "created timestamp not null default now(), primary key(network, hash));"
	_, err := Instance.Execute(query)
	return err
}
