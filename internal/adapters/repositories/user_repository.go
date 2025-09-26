package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"silbackendassessment/internal/core/domain"
	"silbackendassessment/internal/core/ports"
)

type userRepository struct {
	db *bun.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *bun.DB) ports.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user := new(domain.User)
	err := r.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := new(domain.User)
	err := r.db.NewSelect().Model(user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.NewSelect().
		Model(&users).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	return users, err
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	_, err := r.db.NewUpdate().
		Model(user).
		WherePK().
		Exec(ctx)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model((*domain.User)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}