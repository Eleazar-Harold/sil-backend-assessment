package repositories

import (
	"context"
	"database/sql"

	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/core/ports"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type productRepository struct {
	db *bun.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *bun.DB) ports.ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	_, err := r.db.NewInsert().Model(product).Exec(ctx)
	return err
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	product := new(domain.Product)
	err := r.db.NewSelect().
		Model(product).
		Relation("Category").
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return product, nil
}

func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	product := new(domain.Product)
	err := r.db.NewSelect().
		Model(product).
		Relation("Category").
		Where("sku = ?", sku).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return product, nil
}

func (r *productRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	var products []*domain.Product
	err := r.db.NewSelect().
		Model(&products).
		Relation("Category").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	return products, err
}

func (r *productRepository) GetByCategoryID(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*domain.Product, error) {
	var products []*domain.Product
	err := r.db.NewSelect().
		Model(&products).
		Relation("Category").
		Where("category_id = ?", categoryID).
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	return products, err
}

func (r *productRepository) GetActiveProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	var products []*domain.Product
	err := r.db.NewSelect().
		Model(&products).
		Relation("Category").
		Where("is_active = ?", true).
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	return products, err
}

func (r *productRepository) SearchByName(ctx context.Context, name string, limit, offset int) ([]*domain.Product, error) {
	var products []*domain.Product
	err := r.db.NewSelect().
		Model(&products).
		Relation("Category").
		Where("name ILIKE ?", "%"+name+"%").
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	return products, err
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	_, err := r.db.NewUpdate().
		Model(product).
		WherePK().
		Exec(ctx)
	return err
}

func (r *productRepository) UpdateStock(ctx context.Context, id uuid.UUID, stock int) error {
	_, err := r.db.NewUpdate().
		Model((*domain.Product)(nil)).
		Set("stock = ?", stock).
		Set("updated_at = CURRENT_TIMESTAMP").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model((*domain.Product)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}
