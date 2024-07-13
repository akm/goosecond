package goosecond

import (
	"context"
	"database/sql"
)

var EmptyMigrationContext = func(context.Context, *sql.Tx) error { return nil }
