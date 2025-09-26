package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(`
			CREATE TABLE customers (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				first_name VARCHAR(255) NOT NULL,
				last_name VARCHAR(255) NOT NULL,
				email VARCHAR(255) UNIQUE NOT NULL,
				phone VARCHAR(50),
				address TEXT,
				city VARCHAR(100),
				state VARCHAR(100),
				zip_code VARCHAR(20),
				country VARCHAR(100),
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
			);
		`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(`DROP TABLE IF EXISTS customers;`)
		return err
	})
}
