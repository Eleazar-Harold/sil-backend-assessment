package services

import (
	"context"
	"fmt"
	"time"

	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/core/ports"

	"github.com/google/uuid"
)

type categoryService struct {
	categoryRepo ports.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo ports.CategoryRepository) ports.CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) CreateCategory(ctx context.Context, req *domain.CreateCategoryRequest) (*domain.Category, error) {
	// Check if category already exists
	existingCategory, err := s.categoryRepo.GetByName(ctx, req.Name)
	if err == nil && existingCategory != nil {
		return nil, fmt.Errorf("category with name %s already exists", req.Name)
	}

	// Validate parent category if provided
	if req.ParentID != nil {
		parentCategory, err := s.categoryRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get parent category: %w", err)
		}
		if parentCategory == nil {
			return nil, fmt.Errorf("parent category not found")
		}
	}

	// Create new category
	category := &domain.Category{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (s *categoryService) GetCategory(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if category == nil {
		return nil, fmt.Errorf("category not found")
	}

	return category, nil
}

func (s *categoryService) GetCategories(ctx context.Context, limit, offset int) ([]*domain.Category, error) {
	categories, err := s.categoryRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}

func (s *categoryService) GetRootCategories(ctx context.Context) ([]*domain.Category, error) {
	categories, err := s.categoryRepo.GetRootCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get root categories: %w", err)
	}

	return categories, nil
}

func (s *categoryService) GetSubCategories(ctx context.Context, parentID uuid.UUID) ([]*domain.Category, error) {
	categories, err := s.categoryRepo.GetByParentID(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sub categories: %w", err)
	}

	return categories, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id uuid.UUID, req *domain.UpdateCategoryRequest) (*domain.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if category == nil {
		return nil, fmt.Errorf("category not found")
	}

	// Update fields if provided
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Description != nil {
		category.Description = *req.Description
	}
	if req.ParentID != nil {
		// Validate parent category if provided
		if *req.ParentID != uuid.Nil {
			parentCategory, err := s.categoryRepo.GetByID(ctx, *req.ParentID)
			if err != nil {
				return nil, fmt.Errorf("failed to get parent category: %w", err)
			}
			if parentCategory == nil {
				return nil, fmt.Errorf("parent category not found")
			}
		}
		category.ParentID = req.ParentID
	}
	category.UpdatedAt = time.Now()

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get category: %w", err)
	}

	if category == nil {
		return fmt.Errorf("category not found")
	}

	// Check if category has subcategories
	subCategories, err := s.categoryRepo.GetByParentID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check subcategories: %w", err)
	}
	if len(subCategories) > 0 {
		return fmt.Errorf("cannot delete category with subcategories")
	}

	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}
