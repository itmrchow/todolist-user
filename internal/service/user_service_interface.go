package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserService interface {
	RegisterUser(ctx context.Context, req *RegisterReqDTO) error
	LoginUser(ctx context.Context, req *LoginReqDTO) (resp *LoginRespDTO, err error)
}

type RegisterReqDTO struct {
	Name     string
	Email    string
	Password string
}

type LoginReqDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginRespDTO struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Token     string
	ExpiresIn time.Duration
}
