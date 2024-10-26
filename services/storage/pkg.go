package storage

import (
	"github.com/samber/do/v2"
	"libdb.so/e2clicker/services/storage/postgresql"
)

var Package = do.Package(
	postgresql.Package,
)
