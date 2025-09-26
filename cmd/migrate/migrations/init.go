package migrations

import "github.com/uptrace/bun/migrate"

// Migrations is the global migrate.Migrations object used to register migrations
var Migrations = migrate.NewMigrations()