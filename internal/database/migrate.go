package database

import (
	"fmt"
)

// models returns the GORM models to migrate.
func models() []any {
	return []any{
		&Link{},
	}
}

// Migrate runs schema migrations. It is idempotent: repeated runs against an
// already-migrated database are no-ops. It is intended to be invoked once per
// deploy via the `migrate` subcommand, wrapped by the infra repo's init
// container or Job — the application process never mutates the schema at
// startup.
func Migrate(db *DB) error {
	if err := db.AutoMigrate(models()...); err != nil {
		return fmt.Errorf("database migrate: %w", err)
	}
	return nil
}
