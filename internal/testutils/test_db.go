package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// TestDB manages test database connections
type TestDB struct {
	DB *bun.DB
}

// NewTestDB creates a new test database connection
func NewTestDB(t *testing.T) *TestDB {
	// Use test database URL from environment or default
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://testuser:testpass@localhost:5432/test_db?sslmode=disable"
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbURL)))

	// Set connection limits for testing
	sqldb.SetMaxOpenConns(5)
	sqldb.SetMaxIdleConns(2)
	sqldb.SetConnMaxLifetime(0)

	// Test connection
	if err := sqldb.Ping(); err != nil {
		t.Skipf("Skipping test: cannot connect to test database: %v", err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	return &TestDB{DB: db}
}

// Close closes the database connection
func (tdb *TestDB) Close() error {
	return tdb.DB.Close()
}

// Cleanup truncates all tables to ensure clean state between tests
func (tdb *TestDB) Cleanup() error {
	tables := []string{
		"order_items",
		"orders",
		"products",
		"categories",
		"customers",
		"users",
	}

	for _, table := range tables {
		if _, err := tdb.DB.NewTruncateTable().Model((*interface{})(nil)).Table(table).Exec(context.Background()); err != nil {
			// Table might not exist, continue
			log.Printf("Warning: could not truncate table %s: %v", table, err)
		}
	}

	return nil
}

// SetupTestSchema creates the test database schema
func (tdb *TestDB) SetupTestSchema() error {
	// Create tables in dependency order
	schema := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS customers (
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
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS categories (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			parent_id UUID REFERENCES categories(id),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS products (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			sku VARCHAR(100) UNIQUE NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			stock INTEGER NOT NULL DEFAULT 0,
			category_id UUID REFERENCES categories(id),
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			customer_id UUID REFERENCES customers(id),
			order_number VARCHAR(50) UNIQUE NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
			total_amount DECIMAL(10,2) NOT NULL,
			shipping_address TEXT NOT NULL,
			billing_address TEXT NOT NULL,
			notes TEXT,
			order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			shipped_date TIMESTAMP,
			delivered_date TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS order_items (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			order_id UUID REFERENCES orders(id),
			product_id UUID REFERENCES products(id),
			quantity INTEGER NOT NULL,
			unit_price DECIMAL(10,2) NOT NULL,
			total_price DECIMAL(10,2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, sql := range schema {
		if _, err := tdb.DB.Exec(sql); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

// SeedTestData populates the database with test data
func (tdb *TestDB) SeedTestData() error {
	// This would typically insert test data
	// For now, we'll keep it simple and let tests create their own data
	return nil
}
