package services

import (
	"context"
	"fmt"
	"time"

	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/core/ports"

	"github.com/google/uuid"
)

type customerService struct {
	customerRepo ports.CustomerRepository
}

// NewCustomerService creates a new customer service
func NewCustomerService(customerRepo ports.CustomerRepository) ports.CustomerService {
	return &customerService{
		customerRepo: customerRepo,
	}
}

func (s *customerService) CreateCustomer(ctx context.Context, req *domain.CreateCustomerRequest) (*domain.Customer, error) {
	// Check if customer already exists
	existingCustomer, err := s.customerRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingCustomer != nil {
		return nil, fmt.Errorf("customer with email %s already exists", req.Email)
	}

	// Create new customer
	customer := &domain.Customer{
		ID:        uuid.New(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		City:      req.City,
		State:     req.State,
		ZipCode:   req.ZipCode,
		Country:   req.Country,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return customer, nil
}

func (s *customerService) GetCustomer(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}

	return customer, nil
}

func (s *customerService) GetCustomers(ctx context.Context, limit, offset int) ([]*domain.Customer, error) {
	customers, err := s.customerRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}

	return customers, nil
}

func (s *customerService) UpdateCustomer(ctx context.Context, id uuid.UUID, req *domain.UpdateCustomerRequest) (*domain.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	if customer == nil {
		return nil, fmt.Errorf("customer not found")
	}

	// Update fields if provided
	if req.FirstName != nil {
		customer.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		customer.LastName = *req.LastName
	}
	if req.Email != nil {
		customer.Email = *req.Email
	}
	if req.Phone != nil {
		customer.Phone = *req.Phone
	}
	if req.Address != nil {
		customer.Address = *req.Address
	}
	if req.City != nil {
		customer.City = *req.City
	}
	if req.State != nil {
		customer.State = *req.State
	}
	if req.ZipCode != nil {
		customer.ZipCode = *req.ZipCode
	}
	if req.Country != nil {
		customer.Country = *req.Country
	}
	customer.UpdatedAt = time.Now()

	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	return customer, nil
}

func (s *customerService) DeleteCustomer(ctx context.Context, id uuid.UUID) error {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	if customer == nil {
		return fmt.Errorf("customer not found")
	}

	if err := s.customerRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	return nil
}
