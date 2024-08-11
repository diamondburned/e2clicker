// Package postgresqlc contains raw SQL queries for the storage services and
// wrapper functions to use them.
package postgresqlc

import (
	_ "embed"

	"libdb.so/lazymigrate"
)

//go:embed schema.sql
var schemaString string

// Schema is the SQL schema for the storage service.
var Schema = lazymigrate.NewSchema(schemaString)
