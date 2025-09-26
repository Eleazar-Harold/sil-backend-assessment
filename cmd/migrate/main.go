// cmd/migrate/main.go
package main

import (
	"silbackendassessment/cmd/migrate/migrations" // Import migrations package to register migrations

	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"

	"silbackendassessment/internal/config"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	// Get the command from command line arguments
	cmd := parseFlags()

	// Load configuration (prefer file, fallback to environment variables)
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Printf("Failed to load configuration file %s: %v. Falling back to environment variables.", *configPath, err)
		cfg, err = config.LoadFromEnv()
		if err != nil {
			log.Fatalf("Failed to load configuration from environment: %v", err)
		}
	}

	// Connect to the database
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create migrations
	// migrations := migrate.NewMigrations()

	// Debug information
	fmt.Printf("Looking for migrations in: %s\n", cfg.Migrations.Dir)

	// Ensure migration directory exists
	if _, err := os.Stat(cfg.Migrations.Dir); os.IsNotExist(err) {
		if err := os.MkdirAll(cfg.Migrations.Dir, 0755); err != nil {
			log.Fatalf("Failed to create migrations directory: %v", err)
		}
		fmt.Printf("Created migrations directory: %s\n", cfg.Migrations.Dir)

		// Create init.go file
		initFile := fmt.Sprintf("%s/init.go", cfg.Migrations.Dir)
		initContent := `package migrations

import "github.com/uptrace/bun/migrate"

// Migrations is the global migrate.Migrations object used to register migrations
var Migrations = migrate.NewMigrations()
`
		if err := os.WriteFile(initFile, []byte(initContent), 0644); err != nil {
			log.Fatalf("Failed to create migrations init.go file: %v", err)
		}
		fmt.Println("Created migrations/init.go file")

		// If we just created the directory, then there are definitely no migrations yet
		if cmd != "create" && cmd != "init" && cmd != "help" {
			fmt.Println("No migrations exist yet. Create one with:")
			fmt.Printf("  go run cmd/migrate/main.go create your_migration_name\n")
			return
		}
	}

	// Register migrations directory
	// path := os.DirFS(cfg.MigrationDir)
	// if err := migrations.Discover(path); err != nil {
	// 	log.Printf("Warning: Failed to discover migrations: %v", err)
	// 	// Continue execution anyway - this might just mean no migrations exist yet
	// }

	// Create migrator
	migrator := migrate.NewMigrator(db, migrations.Migrations)

	// Debug output - print found migrations using MigrationsWithStatus
	// Only try to get migration status if the command is not init
	if cmd != "init" {
		ms, err := migrator.MigrationsWithStatus(context.Background())
		if err != nil {
			fmt.Printf("Could not get migrations: %v\n", err)
			// This likely means migration tables don't exist yet, suggest initialization
			if cmd != "create" && cmd != "help" {
				fmt.Println("Have you initialized migrations? Try running:")
				fmt.Println("  go run cmd/migrate/main.go init")
			}
		} else {
			fmt.Printf("Found migrations: %d\n", len(ms))
			if len(ms) > 0 {
				fmt.Printf("Migration names: %s\n", formatMigrations(ms))
			} else if cmd != "create" && cmd != "help" {
				fmt.Println("No migrations found. Create one with:")
				fmt.Printf("  go run cmd/migrate/main.go create your_migration_name\n")
			}
		}
	}

	// Execute command
	ctx := context.Background()

	switch cmd {
	case "create":
		name := flag.Arg(1)
		if name == "" {
			log.Fatal("Migration name is required")
		}
		if err := createMigration(cfg.Migrations.Dir, name); err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}

		// After creating a migration, remind user to edit it and run migrations
		fmt.Println("\nReminders:")
		fmt.Println("1. Edit the migration file to implement your schema changes")
		fmt.Println("2. Initialize migrations if you haven't already:")
		fmt.Println("   go run cmd/migrate/main.go init")
		fmt.Println("3. Run the migration:")
		fmt.Println("   go run cmd/migrate/main.go up")

	case "init":
		if err := migrator.Init(ctx); err != nil {
			log.Fatalf("Failed to initialize migrations: %v", err)
		}
		fmt.Println("Migrations initialized")
	case "up":
		if err := migrator.Lock(ctx); err != nil {
			log.Fatalf("Failed to acquire migration lock: %v", err)
		}
		defer migrator.Unlock(ctx)

		group, err := migrator.Migrate(ctx)
		if err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		if group.ID == 0 {
			fmt.Println("No migrations to run")
		} else {
			fmt.Printf("Migrated to %s\n", group)
		}
	case "down":
		if err := migrator.Lock(ctx); err != nil {
			log.Fatalf("Failed to acquire migration lock: %v", err)
		}
		defer migrator.Unlock(ctx)

		group, err := migrator.Rollback(ctx)
		if err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		if group.ID == 0 {
			fmt.Println("No migrations to roll back")
		} else {
			fmt.Printf("Rolled back %s\n", group)
		}
	case "status":
		ms, err := migrator.MigrationsWithStatus(ctx)
		if err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
		fmt.Printf("Migrations: %s\n", formatMigrations(ms))
		fmt.Printf("Unapplied migrations: %s\n", formatMigrations(ms.Unapplied()))
		fmt.Printf("Last migration group: %s\n", ms.LastGroup())
	case "mark_applied":
		migrationName := flag.Arg(1)
		if migrationName == "" {
			log.Fatal("Migration name is required")
		}

		// Find the migration by name in the list of available migrations
		ms, err := migrator.MigrationsWithStatus(ctx)
		if err != nil {
			log.Fatalf("Failed to get migrations: %v", err)
		}

		var found bool
		for _, m := range ms {
			// Check if migration name contains the provided string
			if strings.Contains(m.String(), migrationName) {
				// Mark this specific migration as applied
				if err := migrator.MarkApplied(ctx, &m); err != nil {
					log.Fatalf("Failed to mark migration as applied: %v", err)
				}
				fmt.Printf("Marked migration %s as applied\n", m.String())
				found = true
				break
			}
		}

		if !found {
			log.Fatalf("Migration with name containing '%s' not found", migrationName)
		}

	case "mark_all_applied":
		// Get all migrations
		ms, err := migrator.MigrationsWithStatus(ctx)
		if err != nil {
			log.Fatalf("Failed to get migrations: %v", err)
		}

		if len(ms) == 0 {
			fmt.Println("No migrations found to mark as applied")
			return
		}

		// Get unapplied migrations
		unapplied := ms.Unapplied()
		if len(unapplied) == 0 {
			fmt.Println("All migrations are already marked as applied")
			return
		}

		// Mark each unapplied migration as applied
		count := 0
		for _, m := range unapplied {
			if err := migrator.MarkApplied(ctx, &m); err != nil {
				log.Fatalf("Failed to mark migration %s as applied: %v", m.String(), err)
			}
			fmt.Printf("Marked migration %s as applied\n", m.String())
			count++
		}

		fmt.Printf("Successfully marked %d migrations as applied\n", count)
	case "reset":
		if err := migrator.Lock(ctx); err != nil {
			log.Fatalf("Failed to acquire migration lock: %v", err)
		}
		defer migrator.Unlock(ctx)

		if err := migrator.Reset(ctx); err != nil {
			log.Fatalf("Failed to reset migrations: %v", err)
		}
		fmt.Println("All migrations have been rolled back")
	case "unlock":
		// Use direct SQL query to bypass Bun's safety check
		_, err := db.ExecContext(ctx, "DELETE FROM bun_migration_locks")
		if err != nil {
			log.Fatalf("Failed to release migration locks: %v", err)
		}
		fmt.Println("Migration locks released successfully")
	case "help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printHelp()
		os.Exit(1)
	}
}

// parseFlags parses command line flags and returns the command
func parseFlags() string {
	args := flag.Args()
	if len(args) == 0 {
		printHelp()
		os.Exit(1)
	}
	return args[0]
}

// connectDB establishes a connection to the database
func connectDB(cfg *config.Config) (*bun.DB, error) {
	dsn := cfg.GetDSN()
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Set connection limits
	sqldb.SetMaxOpenConns(10)
	sqldb.SetMaxIdleConns(5)
	sqldb.SetConnMaxLifetime(time.Hour)

	// Check connection
	if err := sqldb.Ping(); err != nil {
		return nil, err
	}

	// Create bun.DB instance
	db := bun.NewDB(sqldb, pgdialect.New())

	return db, nil
}

// createMigration creates a new migration file
func createMigration(dir, name string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	name = strings.ReplaceAll(name, " ", "_")
	fileName := fmt.Sprintf("%s/%s_%s.go", dir, time.Now().UTC().Format("20060102150405"), name)

	content := fmt.Sprintf(`package migrations

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

func init() {
	// Use the global Migrations variable defined in init.go
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] %s... ")
		// TODO: add your migration code here
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] %s... ")
		// TODO: add your rollback code here
		return nil
	})
}
`, name, name)

	// Check if the file already exists
	if _, err := os.Stat(fileName); err == nil {
		return errors.New("migration file already exists")
	} else if !os.IsNotExist(err) {
		return err
	}

	// Create and write to the file
	if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
		return err
	}

	// Create init.go if it doesn't exist
	initFile := fmt.Sprintf("%s/init.go", dir)
	if _, err := os.Stat(initFile); os.IsNotExist(err) {
		initContent := `package migrations

import "github.com/uptrace/bun/migrate"

// Migrations is the global migrate.Migrations object used to register migrations
var Migrations = migrate.NewMigrations()
`
		if err := os.WriteFile(initFile, []byte(initContent), 0644); err != nil {
			return fmt.Errorf("failed to create migrations init.go file: %w", err)
		}
		fmt.Println("Created migrations/init.go file")
	}

	fmt.Printf("Created migration %s\n", fileName)
	return nil
}

// formatMigrations formats a list of migrations for display
func formatMigrations(ms migrate.MigrationSlice) string {
	if len(ms) == 0 {
		return "none"
	}
	names := make([]string, len(ms))
	for i, m := range ms {
		names[i] = m.String()
	}
	return strings.Join(names, ", ")
}

// printHelp prints usage information
func printHelp() {
	fmt.Println("Usage: migrate [command] [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  init                Initialize migration tables")
	fmt.Println("  create [name]       Create a new migration")
	fmt.Println("  up                  Apply all pending migrations")
	fmt.Println("  down                Rollback the last migration group")
	fmt.Println("  reset               Rollback all migrations")
	fmt.Println("  status              Show migration status")
	fmt.Println("  mark_applied [name] Mark a migration as applied without running it")
	fmt.Println("  mark_all_applied    Mark all discovered migrations as applied without running them")
	fmt.Println("  unlock              Release any existing migration locks")
	fmt.Println("  help                Show this help message")
	fmt.Println("\nEnvironment variables:")
	fmt.Println("  MIGRATION_DIR       Directory where migrations are stored (default: cmd/migrate/migrations)")
}
