package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/testutils"

	"github.com/google/uuid"
)

func TestCategoryService_CreateCategory(t *testing.T) {
	mockRepo := testutils.NewMockCategoryRepository()
	service := NewCategoryService(mockRepo)
	ctx := context.Background()

	t.Run("Create category successfully", func(t *testing.T) {
		req := &domain.CreateCategoryRequest{
			Name:        "Electronics",
			Description: "Electronic products",
		}

		category, err := service.CreateCategory(ctx, req)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if category == nil {
			t.Error("Expected category to be created")
		}

		if category.Name != "Electronics" {
			t.Errorf("Expected Name to be 'Electronics', got: %s", category.Name)
		}

		if category.Description != "Electronic products" {
			t.Errorf("Expected Description to be 'Electronic products', got: %s", category.Description)
		}

		if category.ID == uuid.Nil {
			t.Error("Expected category ID to be set")
		}

		if category.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}

		if category.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set")
		}
	})

	t.Run("Create category with parent", func(t *testing.T) {
		// Set up parent category
		parentID := uuid.New()
		parentCategory := &domain.Category{
			ID:        parentID,
			Name:      "Parent Category",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.Categories[parentID] = parentCategory

		req := &domain.CreateCategoryRequest{
			Name:        "Child Category",
			Description: "A child category",
			ParentID:    &parentID,
		}

		category, err := service.CreateCategory(ctx, req)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if category == nil {
			t.Error("Expected category to be created")
		}

		if category.ParentID == nil || *category.ParentID != parentID {
			t.Errorf("Expected ParentID to be %s, got: %v", parentID, category.ParentID)
		}
	})

	t.Run("Create category with existing name", func(t *testing.T) {
		// Set up existing category
		existingCategory := &domain.Category{
			ID:        uuid.New(),
			Name:      "Existing Category",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.Categories[existingCategory.ID] = existingCategory

		req := &domain.CreateCategoryRequest{
			Name:        "Existing Category",
			Description: "A category that already exists",
		}

		category, err := service.CreateCategory(ctx, req)

		if err == nil {
			t.Error("Expected error for existing name")
		}

		if category != nil {
			t.Error("Expected category to be nil")
		}

		if err == nil || err.Error() != "category with name Existing Category already exists" {
			t.Errorf("Expected specific error message, got: %v", err)
		}
	})

	t.Run("Create category with non-existent parent", func(t *testing.T) {
		parentID := uuid.New()

		req := &domain.CreateCategoryRequest{
			Name:        "Child Category 2",
			Description: "A child category",
			ParentID:    &parentID,
		}

		category, err := service.CreateCategory(ctx, req)

		if err == nil {
			t.Error("Expected error for non-existent parent")
		}

		if category != nil {
			t.Error("Expected category to be nil")
		}

		if !errors.Is(err, testutils.ErrCategoryNotFound) {
			t.Errorf("Expected ErrCategoryNotFound, got: %v", err)
		}
	})

	t.Run("Repository error during creation", func(t *testing.T) {
		mockRepo.CreateError = errors.New("database error")

		req := &domain.CreateCategoryRequest{
			Name:        "Test Category",
			Description: "A test category",
		}

		category, err := service.CreateCategory(ctx, req)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if category != nil {
			t.Error("Expected category to be nil")
		}
	})
}

func TestCategoryService_GetCategory(t *testing.T) {
	mockRepo := testutils.NewMockCategoryRepository()
	service := NewCategoryService(mockRepo)
	ctx := context.Background()

	t.Run("Get existing category", func(t *testing.T) {
		categoryID := uuid.New()
		expectedCategory := &domain.Category{
			ID:          categoryID,
			Name:        "Electronics",
			Description: "Electronic products",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		mockRepo.Categories[categoryID] = expectedCategory

		category, err := service.GetCategory(ctx, categoryID)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if category == nil {
			t.Error("Expected category to be returned")
		}

		if category.ID != categoryID {
			t.Errorf("Expected category ID %s, got: %s", categoryID, category.ID)
		}
	})

	t.Run("Get non-existent category", func(t *testing.T) {
		categoryID := uuid.New()

		category, err := service.GetCategory(ctx, categoryID)

		if err == nil {
			t.Error("Expected error for non-existent category")
		}

		if category != nil {
			t.Error("Expected category to be nil")
		}

		if !errors.Is(err, testutils.ErrCategoryNotFound) {
			t.Errorf("Expected ErrCategoryNotFound, got: %v", err)
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		categoryID := uuid.New()
		mockRepo.GetByIDError = errors.New("database error")

		category, err := service.GetCategory(ctx, categoryID)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if category != nil {
			t.Error("Expected category to be nil")
		}
	})
}

func TestCategoryService_GetCategories(t *testing.T) {
	mockRepo := testutils.NewMockCategoryRepository()
	service := NewCategoryService(mockRepo)
	ctx := context.Background()

	t.Run("Get categories successfully", func(t *testing.T) {
		categories := []*domain.Category{
			{
				ID:        uuid.New(),
				Name:      "Category 1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New(),
				Name:      "Category 2",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		mockRepo.AllCategories = categories

		result, err := service.GetCategories(ctx, 10, 0)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("Expected 2 categories, got: %d", len(result))
		}
	})

	t.Run("Repository error", func(t *testing.T) {
		mockRepo.GetAllError = errors.New("database error")

		categories, err := service.GetCategories(ctx, 10, 0)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if categories != nil {
			t.Error("Expected categories to be nil")
		}
	})
}

func TestCategoryService_UpdateCategory(t *testing.T) {
	mockRepo := testutils.NewMockCategoryRepository()
	service := NewCategoryService(mockRepo)
	ctx := context.Background()

	t.Run("Update category successfully", func(t *testing.T) {
		categoryID := uuid.New()
		existingCategory := &domain.Category{
			ID:          categoryID,
			Name:        "Original Category",
			Description: "Original description",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		mockRepo.Categories[categoryID] = existingCategory

		newName := "Updated Category"
		newDescription := "Updated description"
		req := &domain.UpdateCategoryRequest{
			Name:        &newName,
			Description: &newDescription,
		}

		category, err := service.UpdateCategory(ctx, categoryID, req)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if category == nil {
			t.Error("Expected category to be returned")
		}

		if category.Name != "Updated Category" {
			t.Errorf("Expected Name to be 'Updated Category', got: %s", category.Name)
		}

		if category.Description != "Updated description" {
			t.Errorf("Expected Description to be 'Updated description', got: %s", category.Description)
		}
	})

	t.Run("Update non-existent category", func(t *testing.T) {
		categoryID := uuid.New()
		newName := "Updated Category"
		req := &domain.UpdateCategoryRequest{
			Name: &newName,
		}

		category, err := service.UpdateCategory(ctx, categoryID, req)

		if err == nil {
			t.Error("Expected error for non-existent category")
		}

		if category != nil {
			t.Error("Expected category to be nil")
		}

		if !errors.Is(err, testutils.ErrCategoryNotFound) {
			t.Errorf("Expected ErrCategoryNotFound, got: %v", err)
		}
	})

	t.Run("Repository error during update", func(t *testing.T) {
		categoryID := uuid.New()
		existingCategory := &domain.Category{
			ID:        categoryID,
			Name:      "Original Category",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.Categories[categoryID] = existingCategory
		mockRepo.UpdateError = errors.New("database error")

		newName := "Updated Category"
		req := &domain.UpdateCategoryRequest{
			Name: &newName,
		}

		category, err := service.UpdateCategory(ctx, categoryID, req)

		if err == nil {
			t.Error("Expected error from repository")
		}

		if category != nil {
			t.Error("Expected category to be nil")
		}
	})
}

func TestCategoryService_DeleteCategory(t *testing.T) {
	mockRepo := testutils.NewMockCategoryRepository()
	service := NewCategoryService(mockRepo)
	ctx := context.Background()

	t.Run("Delete category successfully", func(t *testing.T) {
		categoryID := uuid.New()
		existingCategory := &domain.Category{
			ID:        categoryID,
			Name:      "Test Category",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.Categories[categoryID] = existingCategory

		err := service.DeleteCategory(ctx, categoryID)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Verify category was deleted
		if _, exists := mockRepo.Categories[categoryID]; exists {
			t.Error("Expected category to be deleted")
		}
	})

	t.Run("Delete non-existent category", func(t *testing.T) {
		categoryID := uuid.New()

		err := service.DeleteCategory(ctx, categoryID)

		if err == nil {
			t.Error("Expected error for non-existent category")
		}

		if !errors.Is(err, testutils.ErrCategoryNotFound) {
			t.Errorf("Expected ErrCategoryNotFound, got: %v", err)
		}
	})

	t.Run("Repository error during deletion", func(t *testing.T) {
		categoryID := uuid.New()
		existingCategory := &domain.Category{
			ID:        categoryID,
			Name:      "Test Category",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.Categories[categoryID] = existingCategory
		mockRepo.DeleteError = errors.New("database error")

		err := service.DeleteCategory(ctx, categoryID)

		if err == nil {
			t.Error("Expected error from repository")
		}
	})
}
