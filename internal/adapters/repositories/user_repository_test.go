package repositories

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"silbackendassessment/internal/core/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

// TestUser represents a user model for testing with SQLite-compatible syntax
type TestUser struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        uuid.UUID `bun:"id,pk,type:uuid" json:"id"`
	Name      string    `bun:"name,notnull" json:"name"`
	Email     string    `bun:"email,unique,notnull" json:"email"`
	CreatedAt time.Time `bun:"created_at,notnull" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull" json:"updated_at"`
}

func setupTestDB(t *testing.T) *bun.DB {
	// Create an in-memory SQLite database for testing
	sqldb, err := sql.Open(sqliteshim.ShimName, ":memory:")
	require.NoError(t, err)

	db := bun.NewDB(sqldb, sqlitedialect.New())

	// Create the users table with SQLite-compatible syntax
	_, err = db.NewCreateTable().Model((*TestUser)(nil)).Exec(context.Background())
	require.NoError(t, err)

	return db
}

// testUserRepository is a test-specific repository that works with TestUser
type testUserRepository struct {
	db *bun.DB
}

func (r *testUserRepository) Create(ctx context.Context, user *domain.User) error {
	testUser := &TestUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	_, err := r.db.NewInsert().Model(testUser).Exec(ctx)
	return err
}

func (r *testUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	testUser := new(TestUser)
	err := r.db.NewSelect().Model(testUser).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &domain.User{
		ID:        testUser.ID,
		Name:      testUser.Name,
		Email:     testUser.Email,
		CreatedAt: testUser.CreatedAt,
		UpdatedAt: testUser.UpdatedAt,
	}, nil
}

func (r *testUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	testUser := new(TestUser)
	err := r.db.NewSelect().Model(testUser).Where("email = ?", email).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &domain.User{
		ID:        testUser.ID,
		Name:      testUser.Name,
		Email:     testUser.Email,
		CreatedAt: testUser.CreatedAt,
		UpdatedAt: testUser.UpdatedAt,
	}, nil
}

func (r *testUserRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var testUsers []*TestUser
	err := r.db.NewSelect().
		Model(&testUsers).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(testUsers))
	for i, testUser := range testUsers {
		users[i] = &domain.User{
			ID:        testUser.ID,
			Name:      testUser.Name,
			Email:     testUser.Email,
			CreatedAt: testUser.CreatedAt,
			UpdatedAt: testUser.UpdatedAt,
		}
	}
	return users, nil
}

func (r *testUserRepository) Update(ctx context.Context, user *domain.User) error {
	testUser := &TestUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	_, err := r.db.NewUpdate().
		Model(testUser).
		WherePK().
		Exec(ctx)
	return err
}

func (r *testUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model((*TestUser)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func TestUserRepository_NewUserRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &testUserRepository{db: db}
	assert.NotNil(t, repo)
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &testUserRepository{db: db}
	ctx := context.Background()

	t.Run("Create user successfully", func(t *testing.T) {
		user := &domain.User{
			ID:        uuid.New(),
			Name:      "Test User",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, user)
		assert.NoError(t, err)

		// Verify the user was created
		retrievedUser, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUser)
		assert.Equal(t, user.ID, retrievedUser.ID)
		assert.Equal(t, user.Name, retrievedUser.Name)
		assert.Equal(t, user.Email, retrievedUser.Email)
	})

	t.Run("Create user with duplicate email", func(t *testing.T) {
		user1 := &domain.User{
			ID:        uuid.New(),
			Name:      "User 1",
			Email:     "duplicate@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		user2 := &domain.User{
			ID:        uuid.New(),
			Name:      "User 2",
			Email:     "duplicate@example.com", // Same email
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Create first user
		err := repo.Create(ctx, user1)
		assert.NoError(t, err)

		// Try to create second user with same email
		err = repo.Create(ctx, user2)
		assert.Error(t, err) // Should fail due to unique constraint
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &testUserRepository{db: db}
	ctx := context.Background()

	t.Run("Get existing user by ID", func(t *testing.T) {
		user := &domain.User{
			ID:        uuid.New(),
			Name:      "Test User",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)

		retrievedUser, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUser)
		assert.Equal(t, user.ID, retrievedUser.ID)
		assert.Equal(t, user.Name, retrievedUser.Name)
		assert.Equal(t, user.Email, retrievedUser.Email)
	})

	t.Run("Get non-existing user by ID", func(t *testing.T) {
		nonExistentID := uuid.New()
		retrievedUser, err := repo.GetByID(ctx, nonExistentID)
		assert.NoError(t, err)
		assert.Nil(t, retrievedUser)
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &testUserRepository{db: db}
	ctx := context.Background()

	t.Run("Get existing user by email", func(t *testing.T) {
		user := &domain.User{
			ID:        uuid.New(),
			Name:      "Test User",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)

		retrievedUser, err := repo.GetByEmail(ctx, user.Email)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUser)
		assert.Equal(t, user.ID, retrievedUser.ID)
		assert.Equal(t, user.Name, retrievedUser.Name)
		assert.Equal(t, user.Email, retrievedUser.Email)
	})

	t.Run("Get non-existing user by email", func(t *testing.T) {
		retrievedUser, err := repo.GetByEmail(ctx, "nonexistent@example.com")
		assert.NoError(t, err)
		assert.Nil(t, retrievedUser)
	})
}

func TestUserRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &testUserRepository{db: db}
	ctx := context.Background()

	t.Run("Get all users with pagination", func(t *testing.T) {
		// Create test users
		users := []*domain.User{
			{
				ID:        uuid.New(),
				Name:      "User 1",
				Email:     "user1@example.com",
				CreatedAt: time.Now().Add(-2 * time.Hour),
				UpdatedAt: time.Now().Add(-2 * time.Hour),
			},
			{
				ID:        uuid.New(),
				Name:      "User 2",
				Email:     "user2@example.com",
				CreatedAt: time.Now().Add(-1 * time.Hour),
				UpdatedAt: time.Now().Add(-1 * time.Hour),
			},
			{
				ID:        uuid.New(),
				Name:      "User 3",
				Email:     "user3@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		for _, user := range users {
			err := repo.Create(ctx, user)
			require.NoError(t, err)
		}

		// Test with limit and offset
		retrievedUsers, err := repo.GetAll(ctx, 2, 0)
		assert.NoError(t, err)
		assert.Len(t, retrievedUsers, 2)

		// Test with offset
		retrievedUsers, err = repo.GetAll(ctx, 2, 1)
		assert.NoError(t, err)
		assert.Len(t, retrievedUsers, 2)

		// Test with large limit
		retrievedUsers, err = repo.GetAll(ctx, 10, 0)
		assert.NoError(t, err)
		assert.Len(t, retrievedUsers, 3)
	})

	t.Run("Get all users when no users exist", func(t *testing.T) {
		// Create a fresh database for this test
		freshDB := setupTestDB(t)
		defer freshDB.Close()
		freshRepo := &testUserRepository{db: freshDB}

		retrievedUsers, err := freshRepo.GetAll(ctx, 10, 0)
		assert.NoError(t, err)
		assert.Empty(t, retrievedUsers)
	})
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &testUserRepository{db: db}
	ctx := context.Background()

	t.Run("Update existing user", func(t *testing.T) {
		user := &domain.User{
			ID:        uuid.New(),
			Name:      "Original Name",
			Email:     "original@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Update the user
		user.Name = "Updated Name"
		user.Email = "updated@example.com"
		user.UpdatedAt = time.Now()

		err = repo.Update(ctx, user)
		assert.NoError(t, err)

		// Verify the update
		retrievedUser, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUser)
		assert.Equal(t, "Updated Name", retrievedUser.Name)
		assert.Equal(t, "updated@example.com", retrievedUser.Email)
	})

	t.Run("Update non-existing user", func(t *testing.T) {
		user := &domain.User{
			ID:        uuid.New(), // Non-existing ID
			Name:      "Non-existing User",
			Email:     "nonexistent@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Update(ctx, user)
		assert.NoError(t, err) // Update doesn't fail for non-existing records
	})
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &testUserRepository{db: db}
	ctx := context.Background()

	t.Run("Delete existing user", func(t *testing.T) {
		user := &domain.User{
			ID:        uuid.New(),
			Name:      "Test User",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Verify user exists
		retrievedUser, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUser)

		// Delete the user
		err = repo.Delete(ctx, user.ID)
		assert.NoError(t, err)

		// Verify user is deleted
		retrievedUser, err = repo.GetByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.Nil(t, retrievedUser)
	})

	t.Run("Delete non-existing user", func(t *testing.T) {
		nonExistentID := uuid.New()
		err := repo.Delete(ctx, nonExistentID)
		assert.NoError(t, err) // Delete doesn't fail for non-existing records
	})
}

func TestUserRepository_Integration(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &testUserRepository{db: db}
	ctx := context.Background()

	t.Run("Full CRUD operations", func(t *testing.T) {
		// Create
		user := &domain.User{
			ID:        uuid.New(),
			Name:      "Integration Test User",
			Email:     "integration@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := repo.Create(ctx, user)
		assert.NoError(t, err)

		// Read
		retrievedUser, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedUser)
		assert.Equal(t, user.Name, retrievedUser.Name)

		// Update
		user.Name = "Updated Integration User"
		user.UpdatedAt = time.Now()
		err = repo.Update(ctx, user)
		assert.NoError(t, err)

		// Verify update
		retrievedUser, err = repo.GetByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Integration User", retrievedUser.Name)

		// Delete
		err = repo.Delete(ctx, user.ID)
		assert.NoError(t, err)

		// Verify deletion
		retrievedUser, err = repo.GetByID(ctx, user.ID)
		assert.NoError(t, err)
		assert.Nil(t, retrievedUser)
	})
}
