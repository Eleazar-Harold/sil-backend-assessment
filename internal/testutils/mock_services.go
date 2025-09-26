package testutils

import (
	"context"
	"errors"
	"time"

	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/core/ports"

	"github.com/google/uuid"
)

// Define error variables
var (
	ErrUserNotFound            = errors.New("user not found")
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrCustomerNotFound        = errors.New("customer not found")
	ErrProductNotFound         = errors.New("product not found")
	ErrCategoryNotFound        = errors.New("category not found")
	ErrOrderNotFound           = errors.New("order not found")
	ErrOrderItemNotFound       = errors.New("order item not found")
	ErrInvalidNotificationType = errors.New("invalid notification type")
)

// MockUserService implements ports.UserService for testing
type MockUserService struct {
	Users map[uuid.UUID]*domain.User
}

func NewMockUserService() *MockUserService {
	return &MockUserService{
		Users: make(map[uuid.UUID]*domain.User),
	}
}

func (m *MockUserService) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	user := &domain.User{
		ID:        uuid.New(),
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.Users[user.ID] = user
	return user, nil
}

func (m *MockUserService) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if user, exists := m.Users[id]; exists {
		return user, nil
	}
	return nil, ErrUserNotFound
}

func (m *MockUserService) GetUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(m.Users))
	for _, user := range m.Users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserService) UpdateUser(ctx context.Context, id uuid.UUID, req *domain.UpdateUserRequest) (*domain.User, error) {
	user, exists := m.Users[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	user.UpdatedAt = time.Now()

	return user, nil
}

func (m *MockUserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if _, exists := m.Users[id]; !exists {
		return ErrUserNotFound
	}
	delete(m.Users, id)
	return nil
}

func (m *MockUserService) SearchUsers(ctx context.Context, query string, limit, offset int) ([]*domain.User, error) {
	users := make([]*domain.User, 0)
	for _, user := range m.Users {
		if user.Name == query || user.Email == query {
			users = append(users, user)
		}
	}
	return users, nil
}

// MockCustomerService implements ports.CustomerService for testing
type MockCustomerService struct {
	Customers map[uuid.UUID]*domain.Customer
}

func NewMockCustomerService() *MockCustomerService {
	return &MockCustomerService{
		Customers: make(map[uuid.UUID]*domain.Customer),
	}
}

func (m *MockCustomerService) CreateCustomer(ctx context.Context, req *domain.CreateCustomerRequest) (*domain.Customer, error) {
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
	m.Customers[customer.ID] = customer
	return customer, nil
}

func (m *MockCustomerService) GetCustomer(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	if customer, exists := m.Customers[id]; exists {
		return customer, nil
	}
	return nil, ErrCustomerNotFound
}

func (m *MockCustomerService) GetCustomers(ctx context.Context, limit, offset int) ([]*domain.Customer, error) {
	customers := make([]*domain.Customer, 0, len(m.Customers))
	for _, customer := range m.Customers {
		customers = append(customers, customer)
	}
	return customers, nil
}

func (m *MockCustomerService) UpdateCustomer(ctx context.Context, id uuid.UUID, req *domain.UpdateCustomerRequest) (*domain.Customer, error) {
	customer, exists := m.Customers[id]
	if !exists {
		return nil, ErrCustomerNotFound
	}

	if req.FirstName != nil {
		customer.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		customer.LastName = *req.LastName
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

	return customer, nil
}

func (m *MockCustomerService) DeleteCustomer(ctx context.Context, id uuid.UUID) error {
	if _, exists := m.Customers[id]; !exists {
		return ErrCustomerNotFound
	}
	delete(m.Customers, id)
	return nil
}

func (m *MockCustomerService) SearchCustomers(ctx context.Context, query string, limit, offset int) ([]*domain.Customer, error) {
	customers := make([]*domain.Customer, 0)
	for _, customer := range m.Customers {
		if customer.FirstName == query || customer.LastName == query || customer.Email == query {
			customers = append(customers, customer)
		}
	}
	return customers, nil
}

// MockNotificationService implements ports.NotificationService for testing
type MockNotificationService struct {
	SentEmails []EmailNotification
	SentSMS    []SMSNotification
}

type EmailNotification struct {
	To       string
	Subject  string
	Body     string
	HTMLBody string
}

type SMSNotification struct {
	PhoneNumber string
	Message     string
}

func NewMockNotificationService() *MockNotificationService {
	return &MockNotificationService{
		SentEmails: make([]EmailNotification, 0),
		SentSMS:    make([]SMSNotification, 0),
	}
}

func (m *MockNotificationService) SendEmail(ctx context.Context, to, subject, body string, htmlBody ...string) error {
	email := EmailNotification{
		To:      to,
		Subject: subject,
		Body:    body,
	}
	if len(htmlBody) > 0 {
		email.HTMLBody = htmlBody[0]
	}
	m.SentEmails = append(m.SentEmails, email)
	return nil
}

func (m *MockNotificationService) SendSMS(ctx context.Context, phoneNumber, message string) error {
	sms := SMSNotification{
		PhoneNumber: phoneNumber,
		Message:     message,
	}
	m.SentSMS = append(m.SentSMS, sms)
	return nil
}

func (m *MockNotificationService) SendNotification(ctx context.Context, req *ports.NotificationRequest) error {
	switch req.Type {
	case ports.EmailNotification:
		return m.SendEmail(ctx, req.To, req.Subject, req.Body, req.HTMLBody)
	case ports.SMSNotification:
		return m.SendSMS(ctx, req.PhoneNumber, req.Message)
	default:
		return ErrInvalidNotificationType
	}
}

func (m *MockNotificationService) SendBulkEmail(ctx context.Context, recipients []string, subject, body string, htmlBody ...string) error {
	for _, recipient := range recipients {
		if err := m.SendEmail(ctx, recipient, subject, body, htmlBody...); err != nil {
			return err
		}
	}
	return nil
}

func (m *MockNotificationService) SendBulkSMS(ctx context.Context, phoneNumbers []string, message string) error {
	for _, phoneNumber := range phoneNumbers {
		if err := m.SendSMS(ctx, phoneNumber, message); err != nil {
			return err
		}
	}
	return nil
}

func (m *MockNotificationService) ValidateEmail(email string) bool {
	return len(email) > 0 && email != "invalid@" && email != "invalid-email"
}

func (m *MockNotificationService) ValidatePhoneNumber(phoneNumber string) bool {
	return len(phoneNumber) > 0 && phoneNumber != "invalid"
}

func (m *MockNotificationService) SendOrderConfirmationEmail(ctx context.Context, customerEmail, customerName, orderNumber string, orderItems []domain.OrderItem) error {
	return m.SendEmail(ctx, customerEmail, "Order Confirmation", "Your order has been confirmed")
}

func (m *MockNotificationService) SendOrderConfirmationSMS(ctx context.Context, phoneNumber, customerName, orderNumber string, totalAmount float64) error {
	return m.SendSMS(ctx, phoneNumber, "Your order has been confirmed")
}

func (m *MockNotificationService) SendOrderStatusUpdateSMS(ctx context.Context, phoneNumber, customerName, orderNumber, status string) error {
	return m.SendSMS(ctx, phoneNumber, "Order status updated")
}

// MockEmailClient implements ports.EmailClient for testing
type MockEmailClient struct {
	SentEmails []EmailNotification
}

func (m *MockEmailClient) SendEmail(ctx context.Context, to, subject, body string, htmlBody ...string) error {
	email := EmailNotification{
		To:      to,
		Subject: subject,
		Body:    body,
	}
	if len(htmlBody) > 0 {
		email.HTMLBody = htmlBody[0]
	}
	m.SentEmails = append(m.SentEmails, email)
	return nil
}

func (m *MockEmailClient) SendBulkEmail(ctx context.Context, recipients []string, subject, body string, htmlBody ...string) error {
	for _, recipient := range recipients {
		if err := m.SendEmail(ctx, recipient, subject, body, htmlBody...); err != nil {
			return err
		}
	}
	return nil
}

func (m *MockEmailClient) ValidateEmail(email string) bool {
	return len(email) > 0 && email != "invalid@" && email != "invalid-email"
}

// MockSMSClient implements ports.SMSClient for testing
type MockSMSClient struct {
	SentSMS []SMSNotification
}

func (m *MockSMSClient) SendSMS(ctx context.Context, phoneNumber, message string) error {
	sms := SMSNotification{
		PhoneNumber: phoneNumber,
		Message:     message,
	}
	m.SentSMS = append(m.SentSMS, sms)
	return nil
}

func (m *MockSMSClient) SendBulkSMS(ctx context.Context, phoneNumbers []string, message string) error {
	for _, phoneNumber := range phoneNumbers {
		if err := m.SendSMS(ctx, phoneNumber, message); err != nil {
			return err
		}
	}
	return nil
}

func (m *MockSMSClient) ValidatePhoneNumber(phoneNumber string) bool {
	return len(phoneNumber) > 0 && phoneNumber != "invalid"
}

// MockUserRepository implements ports.UserRepository for testing
type MockUserRepository struct {
	Users           map[uuid.UUID]*domain.User
	UsersByEmail    map[string]*domain.User
	AllUsers        []*domain.User
	CreateError     error
	GetByIDError    error
	GetByEmailError error
	GetAllError     error
	UpdateError     error
	DeleteError     error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users:        make(map[uuid.UUID]*domain.User),
		UsersByEmail: make(map[string]*domain.User),
		AllUsers:     make([]*domain.User, 0),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	m.Users[user.ID] = user
	m.UsersByEmail[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if m.GetByIDError != nil {
		return nil, m.GetByIDError
	}
	if user, exists := m.Users[id]; exists {
		return user, nil
	}
	return nil, ErrUserNotFound
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.GetByEmailError != nil {
		return nil, m.GetByEmailError
	}
	if user, exists := m.UsersByEmail[email]; exists {
		return user, nil
	}
	return nil, ErrUserNotFound
}

func (m *MockUserRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	if m.GetAllError != nil {
		return nil, m.GetAllError
	}
	return m.AllUsers, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	m.Users[user.ID] = user
	m.UsersByEmail[user.Email] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	if user, exists := m.Users[id]; exists {
		delete(m.Users, id)
		delete(m.UsersByEmail, user.Email)
		return nil
	}
	return ErrUserNotFound
}

// MockCustomerRepository implements ports.CustomerRepository for testing
type MockCustomerRepository struct {
	Customers        map[uuid.UUID]*domain.Customer
	CustomersByEmail map[string]*domain.Customer
	AllCustomers     []*domain.Customer
	CreateError      error
	GetByIDError     error
	GetByEmailError  error
	GetAllError      error
	UpdateError      error
	DeleteError      error
}

func NewMockCustomerRepository() *MockCustomerRepository {
	return &MockCustomerRepository{
		Customers:        make(map[uuid.UUID]*domain.Customer),
		CustomersByEmail: make(map[string]*domain.Customer),
		AllCustomers:     make([]*domain.Customer, 0),
	}
}

func (m *MockCustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	m.Customers[customer.ID] = customer
	m.CustomersByEmail[customer.Email] = customer
	return nil
}

func (m *MockCustomerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	if m.GetByIDError != nil {
		return nil, m.GetByIDError
	}
	if customer, exists := m.Customers[id]; exists {
		return customer, nil
	}
	return nil, ErrCustomerNotFound
}

func (m *MockCustomerRepository) GetByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	if m.GetByEmailError != nil {
		return nil, m.GetByEmailError
	}
	if customer, exists := m.CustomersByEmail[email]; exists {
		return customer, nil
	}
	return nil, ErrCustomerNotFound
}

func (m *MockCustomerRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Customer, error) {
	if m.GetAllError != nil {
		return nil, m.GetAllError
	}
	return m.AllCustomers, nil
}

func (m *MockCustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	m.Customers[customer.ID] = customer
	m.CustomersByEmail[customer.Email] = customer
	return nil
}

func (m *MockCustomerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	if customer, exists := m.Customers[id]; exists {
		delete(m.Customers, id)
		delete(m.CustomersByEmail, customer.Email)
		return nil
	}
	return ErrCustomerNotFound
}

// MockProductRepository implements ports.ProductRepository for testing
type MockProductRepository struct {
	Products      map[uuid.UUID]*domain.Product
	ProductsBySKU map[string]*domain.Product
	AllProducts   []*domain.Product
	CreateError   error
	GetByIDError  error
	GetBySKUError error
	GetAllError   error
	UpdateError   error
	DeleteError   error
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		Products:      make(map[uuid.UUID]*domain.Product),
		ProductsBySKU: make(map[string]*domain.Product),
		AllProducts:   make([]*domain.Product, 0),
	}
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	m.Products[product.ID] = product
	m.ProductsBySKU[product.SKU] = product
	return nil
}

func (m *MockProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	if m.GetByIDError != nil {
		return nil, m.GetByIDError
	}
	if product, exists := m.Products[id]; exists {
		return product, nil
	}
	return nil, ErrProductNotFound
}

func (m *MockProductRepository) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	if m.GetBySKUError != nil {
		return nil, m.GetBySKUError
	}
	if product, exists := m.ProductsBySKU[sku]; exists {
		return product, nil
	}
	return nil, ErrProductNotFound
}

func (m *MockProductRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	if m.GetAllError != nil {
		return nil, m.GetAllError
	}
	return m.AllProducts, nil
}

func (m *MockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	m.Products[product.ID] = product
	m.ProductsBySKU[product.SKU] = product
	return nil
}

func (m *MockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	if product, exists := m.Products[id]; exists {
		delete(m.Products, id)
		delete(m.ProductsBySKU, product.SKU)
		return nil
	}
	return ErrProductNotFound
}

func (m *MockProductRepository) GetByCategoryID(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*domain.Product, error) {
	// Simple implementation for testing
	return m.AllProducts, nil
}

func (m *MockProductRepository) GetActiveProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	// Simple implementation for testing
	return m.AllProducts, nil
}

func (m *MockProductRepository) SearchByName(ctx context.Context, name string, limit, offset int) ([]*domain.Product, error) {
	// Simple implementation for testing
	return m.AllProducts, nil
}

func (m *MockProductRepository) UpdateStock(ctx context.Context, id uuid.UUID, stock int) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	if product, exists := m.Products[id]; exists {
		product.Stock = stock
		return nil
	}
	return ErrProductNotFound
}

// MockCategoryRepository implements ports.CategoryRepository for testing
type MockCategoryRepository struct {
	Categories    map[uuid.UUID]*domain.Category
	AllCategories []*domain.Category
	CreateError   error
	GetByIDError  error
	GetAllError   error
	UpdateError   error
	DeleteError   error
}

func NewMockCategoryRepository() *MockCategoryRepository {
	return &MockCategoryRepository{
		Categories:    make(map[uuid.UUID]*domain.Category),
		AllCategories: make([]*domain.Category, 0),
	}
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	m.Categories[category.ID] = category
	return nil
}

func (m *MockCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	if m.GetByIDError != nil {
		return nil, m.GetByIDError
	}
	if category, exists := m.Categories[id]; exists {
		return category, nil
	}
	return nil, ErrCategoryNotFound
}

func (m *MockCategoryRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Category, error) {
	if m.GetAllError != nil {
		return nil, m.GetAllError
	}
	return m.AllCategories, nil
}

func (m *MockCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	m.Categories[category.ID] = category
	return nil
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	if _, exists := m.Categories[id]; exists {
		delete(m.Categories, id)
		return nil
	}
	return ErrCategoryNotFound
}

func (m *MockCategoryRepository) GetByName(ctx context.Context, name string) (*domain.Category, error) {
	// Simple implementation for testing
	for _, category := range m.Categories {
		if category.Name == name {
			return category, nil
		}
	}
	return nil, ErrCategoryNotFound
}

func (m *MockCategoryRepository) GetByParentID(ctx context.Context, parentID uuid.UUID) ([]*domain.Category, error) {
	// Simple implementation for testing
	return m.AllCategories, nil
}

func (m *MockCategoryRepository) GetRootCategories(ctx context.Context) ([]*domain.Category, error) {
	// Simple implementation for testing
	return m.AllCategories, nil
}

// MockOrderRepository implements ports.OrderRepository for testing
type MockOrderRepository struct {
	Orders       map[uuid.UUID]*domain.Order
	AllOrders    []*domain.Order
	CreateError  error
	GetByIDError error
	GetAllError  error
	UpdateError  error
	DeleteError  error
}

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{
		Orders:    make(map[uuid.UUID]*domain.Order),
		AllOrders: make([]*domain.Order, 0),
	}
}

func (m *MockOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	m.Orders[order.ID] = order
	return nil
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	if m.GetByIDError != nil {
		return nil, m.GetByIDError
	}
	if order, exists := m.Orders[id]; exists {
		return order, nil
	}
	return nil, ErrOrderNotFound
}

func (m *MockOrderRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Order, error) {
	if m.GetAllError != nil {
		return nil, m.GetAllError
	}
	return m.AllOrders, nil
}

func (m *MockOrderRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*domain.Order, error) {
	// Simple implementation for testing
	return m.AllOrders, nil
}

func (m *MockOrderRepository) GetByStatus(ctx context.Context, status domain.OrderStatus, limit, offset int) ([]*domain.Order, error) {
	// Simple implementation for testing
	return m.AllOrders, nil
}

func (m *MockOrderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	// Simple implementation for testing
	for _, order := range m.Orders {
		if order.OrderNumber == orderNumber {
			return order, nil
		}
	}
	return nil, ErrOrderNotFound
}

func (m *MockOrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	if order, exists := m.Orders[id]; exists {
		order.Status = status
		return nil
	}
	return ErrOrderNotFound
}

func (m *MockOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	m.Orders[order.ID] = order
	return nil
}

func (m *MockOrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	if _, exists := m.Orders[id]; exists {
		delete(m.Orders, id)
		return nil
	}
	return ErrOrderNotFound
}

// MockOrderItemRepository implements ports.OrderItemRepository for testing
type MockOrderItemRepository struct {
	OrderItems    map[uuid.UUID]*domain.OrderItem
	AllOrderItems []*domain.OrderItem
	CreateError   error
	GetByIDError  error
	GetAllError   error
	UpdateError   error
	DeleteError   error
}

func NewMockOrderItemRepository() *MockOrderItemRepository {
	return &MockOrderItemRepository{
		OrderItems:    make(map[uuid.UUID]*domain.OrderItem),
		AllOrderItems: make([]*domain.OrderItem, 0),
	}
}

func (m *MockOrderItemRepository) Create(ctx context.Context, orderItem *domain.OrderItem) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	m.OrderItems[orderItem.ID] = orderItem
	return nil
}

func (m *MockOrderItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.OrderItem, error) {
	if m.GetByIDError != nil {
		return nil, m.GetByIDError
	}
	if orderItem, exists := m.OrderItems[id]; exists {
		return orderItem, nil
	}
	return nil, ErrOrderItemNotFound
}

func (m *MockOrderItemRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.OrderItem, error) {
	if m.GetAllError != nil {
		return nil, m.GetAllError
	}
	return m.AllOrderItems, nil
}

func (m *MockOrderItemRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderItem, error) {
	// Simple implementation for testing
	return m.AllOrderItems, nil
}

func (m *MockOrderItemRepository) Update(ctx context.Context, orderItem *domain.OrderItem) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	m.OrderItems[orderItem.ID] = orderItem
	return nil
}

func (m *MockOrderItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	if _, exists := m.OrderItems[id]; exists {
		delete(m.OrderItems, id)
		return nil
	}
	return ErrOrderItemNotFound
}

func (m *MockOrderItemRepository) GetByProductID(ctx context.Context, productID uuid.UUID, limit, offset int) ([]*domain.OrderItem, error) {
	// Simple implementation for testing
	return m.AllOrderItems, nil
}

func (m *MockOrderItemRepository) DeleteByOrderID(ctx context.Context, orderID uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	// Simple implementation for testing - remove all items for the order
	for id, item := range m.OrderItems {
		if item.OrderID == orderID {
			delete(m.OrderItems, id)
		}
	}
	return nil
}
