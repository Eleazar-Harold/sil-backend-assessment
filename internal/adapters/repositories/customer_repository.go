package repositories

import (
	"context"
	"database/sql"

	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/core/ports"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type customerRepository struct {
	db *bun.DB
}

// NewCustomerRepository creates a new customer repository
func NewCustomerRepository(db *bun.DB) ports.CustomerRepository {
	return &customerRepository{
		db: db,
	}
}

func (r *customerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	_, err := r.db.NewInsert().Model(customer).Exec(ctx)
	return err
}

func (r *customerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	customer := new(domain.Customer)
	err := r.db.NewSelect().Model(customer).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return customer, nil
}

func (r *customerRepository) GetByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	customer := new(domain.Customer)
	err := r.db.NewSelect().Model(customer).Where("email = ?", email).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return customer, nil
}

func (r *customerRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Customer, error) {
	var customers []*domain.Customer
	err := r.db.NewSelect().
		Model(&customers).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	return customers, err
}

func (r *customerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	_, err := r.db.NewUpdate().
		Model(customer).
		WherePK().
		Exec(ctx)
	return err
}

func (r *customerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model((*domain.Customer)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}
