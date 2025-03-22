package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/itmrchow/todolist-proto/protobuf"
	pb "github.com/itmrchow/todolist-proto/protobuf/user"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/itmrchow/todolist-users/internal/entity"
	mErr "github.com/itmrchow/todolist-users/internal/errors"
	"github.com/itmrchow/todolist-users/internal/repository"
	"github.com/itmrchow/todolist-users/utils"
)

var (
	secretKey = viper.GetString("jwt.secret_key")
	expireAt  = viper.GetInt("jwt.expire_at")
	issuer    = viper.GetString("server_name")
)

type userServiceImpl struct {
	pb.UnimplementedUserServiceServer
	userRepo repository.UsersRepository
}

func NewUserService(userRepo repository.UsersRepository) pb.UserServiceServer {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

func (u *userServiceImpl) Login(ctx context.Context, req *pb.LoginRequest) (resp *pb.LoginResponse, err error) {
	// get user info by email
	user := &entity.User{
		Email:    req.Email,
		Password: req.Password,
	}
	user.HashPassword()

	user, err = u.userRepo.GetByEmailAndPassword(ctx, user.Email, user.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.Unauthenticated, mErr.ErrInvalidLoginInfo)
		}
		log.Error().Err(err).Msg("GetByEmailAndPassword error")
		return nil, status.Error(codes.Internal, mErr.ErrInternalServerError)
	}

	// generate token
	token, err := utils.GenerateToken(user.ID.String(), secretKey, issuer, expireAt)
	if err != nil {
		log.Error().Err(err).Msg("Generate token error")
		return nil, status.Error(codes.Internal, mErr.ErrInternalServerError)
	}

	resp = &pb.LoginResponse{
		Id:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Token:     token,
		ExpiresIn: timestamppb.New(time.Now().Add(time.Duration(expireAt) * time.Hour)),
	}

	return

}

func (u *userServiceImpl) Register(ctx context.Context, req *pb.RegisterRequest) (resp *protobuf.EmptyResponse, err error) {
	// check email is exist?
	isExist, err := u.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Msg("internal server error")
		return nil, status.Error(codes.Internal, mErr.ErrInternalServerError)
	}

	if isExist {
		return nil, status.Error(codes.AlreadyExists, mErr.ErrEmailAlreadyExists)
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
		log.Error().Err(err).Msg("user , insert db error")
		return nil, status.Error(codes.Internal, mErr.ErrInternalServerError)
	}

	return
}
