package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/itmrchow/todolist-users/internal/entity"
	mErr "github.com/itmrchow/todolist-users/internal/errors"
	"github.com/itmrchow/todolist-users/internal/repository"
	"github.com/itmrchow/todolist-users/utils"
)

var _ UserService = &userServiceImpl{}

var (
	secretKey = viper.GetString("jwt.secret_key")
	expireAt  = viper.GetInt("jwt.expire_at")
	issuer    = viper.GetString("server_name")
)

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
		return &mErr.Err500InternalServer
	}

	if isExist {
		return &mErr.Err400EmailAlreadyExists
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
		return &mErr.Err500InternalServer
	}

	return
}

func (u *userServiceImpl) LoginUser(ctx context.Context, req *LoginReqDTO) (resp *LoginRespDTO, err error) {

	// get user info by email
	user := &entity.User{
		Email:    req.Email,
		Password: req.Password,
	}
	user.HashPassword()

	user, err = u.userRepo.GetByEmailAndPassword(ctx, user.Email, user.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &mErr.Err400InvalidLoginInfo
		}
		// TODO: print log
		return nil, &mErr.Err500InternalServer
	}

	// generate token
	token, err := utils.GenerateToken(user.ID.String(), secretKey, issuer, expireAt)
	if err != nil {
		// TODO: print log
		return nil, &mErr.Err500InternalServer
	}

	resp = &LoginRespDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Token:     token,
		ExpiresIn: time.Now().Add(time.Duration(expireAt) * time.Hour),
	}

	return
}
