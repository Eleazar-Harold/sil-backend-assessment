package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(`
			CREATE TABLE categories (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				name VARCHAR(255) UNIQUE NOT NULL,
				description TEXT,
				parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
			);
		`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(`DROP TABLE IF EXISTS categories;`)
		return err
	})
}
