package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(`
			CREATE TABLE orders (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
				order_number VARCHAR(100) UNIQUE NOT NULL,
				status VARCHAR(50) NOT NULL DEFAULT 'pending',
				total_amount DECIMAL(10,2) NOT NULL,
				shipping_address TEXT NOT NULL,
				billing_address TEXT NOT NULL,
				notes TEXT,
				order_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				shipped_date TIMESTAMP,
				delivered_date TIMESTAMP,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
			);
		`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(`DROP TABLE IF EXISTS orders;`)
		return err
	})
}
