package storage

import (
	"go.uber.org/fx"
	"libdb.so/e2clicker/services/storage/postgresql"
)

var Module = fx.Module("storage",
	postgresql.Module,
)
