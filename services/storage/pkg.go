package storage

import (
	"go.uber.org/fx"
	"e2clicker.app/services/storage/postgresql"
)

var Module = fx.Module("storage",
	postgresql.Module,
)
