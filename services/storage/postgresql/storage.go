package postgresql

import "github.com/jackc/pgx/v5/pgxpool"

// Storage is the PostgreSQL-backed storage.
type Storage struct {
	conn *pgxpool.Pool
}

func newStorage(conn *pgxpool.Pool) *Storage {
	return &Storage{conn: conn}
}
