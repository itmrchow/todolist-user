package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/itmrchow/todolist-users/internal/entity"
	mErr "github.com/itmrchow/todolist-users/internal/errors"
	"github.com/itmrchow/todolist-users/internal/repository"
)

var _ UserService = &userServiceImpl{}

type userServiceImpl struct {
	userRepo repository.UsersRepository
}

func NewUserService(userRepo repository.UsersRepository) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

func (u *userServiceImpl) RegisterUser(ctx context.Context, req *RegisterReqDTO) (err error) {

	// check email is exist?
	isExist, err := u.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		// TODO: print log
		return &mErr.InternalServerError{}
	}

	if isExist {
		return &mErr.BadRequestError{Msg: mErr.MsgEmailAlreadyExists}
	}

	// insert db
	user := &entity.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	user.HashPassword()

	err = u.userRepo.Create(ctx, user)

	if err != nil {
		// TODO: print log
		return &mErr.InternalServerError{}
	}

	return
}

func (u *userServiceImpl) LoginUser(ctx context.Context, req *LoginReqDTO) (resp *LoginRespDTO, err error) {

	// get user info by email

	// check password is correct

	// generate token

	// return token

	return nil, nil
}
