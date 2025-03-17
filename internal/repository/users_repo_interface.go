package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/itmrchow/todolist-users/internal/entity"
)

type UsersRepository interface {
	Create(ctx context.Context, user *entity.User) error
	Get(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
