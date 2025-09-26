package repositories

import (
	"context"
	"database/sql"

	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/core/ports"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type categoryRepository struct {
	db *bun.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *bun.DB) ports.CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	_, err := r.db.NewInsert().Model(category).Exec(ctx)
	return err
}

func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	category := new(domain.Category)
	err := r.db.NewSelect().Model(category).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return category, nil
}

func (r *categoryRepository) GetByName(ctx context.Context, name string) (*domain.Category, error) {
	category := new(domain.Category)
	err := r.db.NewSelect().Model(category).Where("name = ?", name).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return category, nil
}

func (r *categoryRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.NewSelect().
		Model(&categories).
		Order("name ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	return categories, err
}

func (r *categoryRepository) GetByParentID(ctx context.Context, parentID uuid.UUID) ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.NewSelect().
		Model(&categories).
		Where("parent_id = ?", parentID).
		Order("name ASC").
		Scan(ctx)
	return categories, err
}

func (r *categoryRepository) GetRootCategories(ctx context.Context) ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.NewSelect().
		Model(&categories).
		Where("parent_id IS NULL").
		Order("name ASC").
		Scan(ctx)
	return categories, err
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	_, err := r.db.NewUpdate().
		Model(category).
		WherePK().
		Exec(ctx)
	return err
}

func (r *categoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model((*domain.Category)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}
