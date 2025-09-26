package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(`
			CREATE TABLE products (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				name VARCHAR(255) NOT NULL,
				description TEXT,
				sku VARCHAR(100) UNIQUE NOT NULL,
				price DECIMAL(10,2) NOT NULL,
				stock INTEGER NOT NULL DEFAULT 0,
				category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
				is_active BOOLEAN NOT NULL DEFAULT true,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
			);
		`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(`DROP TABLE IF EXISTS products;`)
		return err
	})
}
