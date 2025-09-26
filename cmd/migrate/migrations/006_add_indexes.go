package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		// Customer indexes
		_, err := db.Exec(`CREATE INDEX idx_customers_email ON customers(email);`)
		if err != nil {
			return err
		}

		// Category indexes
		_, err = db.Exec(`CREATE INDEX idx_categories_parent_id ON categories(parent_id);`)
		if err != nil {
			return err
		}

		// Product indexes
		_, err = db.Exec(`CREATE INDEX idx_products_sku ON products(sku);`)
		if err != nil {
			return err
		}
		_, err = db.Exec(`CREATE INDEX idx_products_category_id ON products(category_id);`)
		if err != nil {
			return err
		}
		_, err = db.Exec(`CREATE INDEX idx_products_is_active ON products(is_active);`)
		if err != nil {
			return err
		}

		// Order indexes
		_, err = db.Exec(`CREATE INDEX idx_orders_customer_id ON orders(customer_id);`)
		if err != nil {
			return err
		}
		_, err = db.Exec(`CREATE INDEX idx_orders_order_number ON orders(order_number);`)
		if err != nil {
			return err
		}
		_, err = db.Exec(`CREATE INDEX idx_orders_status ON orders(status);`)
		if err != nil {
			return err
		}
		_, err = db.Exec(`CREATE INDEX idx_orders_order_date ON orders(order_date);`)
		if err != nil {
			return err
		}

		// Order item indexes
		_, err = db.Exec(`CREATE INDEX idx_order_items_order_id ON order_items(order_id);`)
		if err != nil {
			return err
		}
		_, err = db.Exec(`CREATE INDEX idx_order_items_product_id ON order_items(product_id);`)
		if err != nil {
			return err
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		// Drop indexes in reverse order
		indexes := []string{
			"idx_order_items_product_id",
			"idx_order_items_order_id",
			"idx_orders_order_date",
			"idx_orders_status",
			"idx_orders_order_number",
			"idx_orders_customer_id",
			"idx_products_is_active",
			"idx_products_category_id",
			"idx_products_sku",
			"idx_categories_parent_id",
			"idx_customers_email",
		}

		for _, index := range indexes {
			_, err := db.Exec(`DROP INDEX IF EXISTS ` + index + `;`)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
