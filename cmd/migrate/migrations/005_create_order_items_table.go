package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(`
			CREATE TABLE order_items (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
				product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
				quantity INTEGER NOT NULL,
				unit_price DECIMAL(10,2) NOT NULL,
				total_price DECIMAL(10,2) NOT NULL,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
			);
		`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.Exec(`DROP TABLE IF EXISTS order_items;`)
		return err
	})
}
