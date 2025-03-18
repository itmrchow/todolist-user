package service

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/itmrchow/todolist-users/internal/entity"
	mErr "github.com/itmrchow/todolist-users/internal/errors"
	"github.com/itmrchow/todolist-users/internal/repository"
)

func TestUserServiceImplSuite(t *testing.T) {
	suite.Run(t, new(UserServiceImplTestSuite))
}

type UserServiceImplTestSuite struct {
	suite.Suite
	userService  UserService
	mockUserRepo *repository.MockUsersRepository
}

func (s *UserServiceImplTestSuite) SetupTest() {

	userRepo := repository.NewMockUsersRepository(s.T())
	s.mockUserRepo = userRepo
	s.userService = NewUserService(userRepo)
}

func (s *UserServiceImplTestSuite) Test_userServiceImpl_RegisterUser() {
	type args struct {
		ctx context.Context
		req *RegisterReqDTO
	}
	tests := []struct {
		name       string
		args       args
		mockFunc   func(m *repository.MockUsersRepository)
		assertFunc func(err error)
	}{
		{
			name: "db error",
			args: args{
				ctx: context.Background(),
				req: &RegisterReqDTO{
					Email:    "db_error@example.com",
					Password: "password",
					Name:     "test",
				},
			},
			mockFunc: func(m *repository.MockUsersRepository) {
				m.EXPECT().ExistsByEmail(context.Background(), "db_error@example.com").Return(false, errors.New("db error"))
			},
			assertFunc: func(err error) {
				s.Assert().Error(err)
				s.Assert().ErrorIs(err, &mErr.Err500InternalServer)
			},
		},
		{
			name: "email exist",
			args: args{
				ctx: context.Background(),
				req: &RegisterReqDTO{
					Email:    "exist@example.com",
					Password: "password",
					Name:     "test",
				},
			},
			mockFunc: func(mock *repository.MockUsersRepository) {
				mock.EXPECT().ExistsByEmail(context.Background(), "exist@example.com").Return(true, nil)
			},
			assertFunc: func(err error) {
				s.Assert().Error(err)
				s.Assert().Equal("email already exists", err.Error())
			},
		},
		{
			name: "insert db fail",
			args: args{
				ctx: context.Background(),
				req: &RegisterReqDTO{
					Email:    "test@example.com",
					Password: "password",
					Name:     "test",
				},
			},
			mockFunc: func(m *repository.MockUsersRepository) {
				m.EXPECT().ExistsByEmail(context.Background(), "test@example.com").Return(false, nil)
				m.EXPECT().Create(context.Background(), mock.Anything).Return(errors.New("db error")).Times(1)
			},
			assertFunc: func(err error) {
				s.Assert().Error(err)
				s.Assert().Equal("internal server error", err.Error())
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: &RegisterReqDTO{
					Email:    "test@example.com",
					Password: "password",
					Name:     "test",
				},
			},
			mockFunc: func(m *repository.MockUsersRepository) {
				m.EXPECT().ExistsByEmail(context.Background(), "test@example.com").Return(false, nil)
				m.EXPECT().Create(context.Background(), mock.MatchedBy(func(user *entity.User) bool {

					passwordBytes := sha512.Sum512([]byte(user.Password))
					password := fmt.Sprintf("%x", passwordBytes)

					return user.Email == "test@example.com" && user.Password != password && user.Name == "test"

				})).Return(nil).Times(1)
			},
			assertFunc: func(err error) {
				s.Assert().NoError(err)
			},
		},
	}
	for _, tt := range tests {

		s.Run(tt.name, func() {

			tt.mockFunc(s.mockUserRepo)
			err := s.userService.RegisterUser(tt.args.ctx, tt.args.req)
			tt.assertFunc(err)
		})
	}
}

func (s *UserServiceImplTestSuite) Test_userServiceImpl_LoginUser() {
	type args struct {
		ctx context.Context
		req *LoginReqDTO
	}
	tests := []struct {
		name       string
		args       args
		mockFunc   func(m *repository.MockUsersRepository)
		assertFunc func(resp *LoginRespDTO, err error)
	}{
		{
			name: "db error",
			args: args{
				ctx: context.Background(),
				req: &LoginReqDTO{
					Email:    "db_error@example.com",
					Password: "password",
				},
			},
			mockFunc: func(m *repository.MockUsersRepository) {
				m.EXPECT().GetByEmailAndPassword(context.Background(), "db_error@example.com", mock.AnythingOfType("string")).Return(nil, errors.New("db error")).Times(1)
			},
			assertFunc: func(resp *LoginRespDTO, err error) {
				s.Assert().ErrorIs(err, &mErr.Err500InternalServer)
			},
		},
		{
			name: "user not found",
			args: args{
				ctx: context.Background(),
				req: &LoginReqDTO{
					Email:    "not_found@example.com",
					Password: "password",
				},
			},
			mockFunc: func(m *repository.MockUsersRepository) {
				m.EXPECT().GetByEmailAndPassword(context.Background(), "not_found@example.com", mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound).Times(1)
			},
			assertFunc: func(resp *LoginRespDTO, err error) {
				s.Assert().ErrorIs(err, &mErr.Err400InvalidLoginInfo)
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: &LoginReqDTO{
					Email:    "success@example.com",
					Password: "password",
				},
			},
			mockFunc: func(m *repository.MockUsersRepository) {
				m.EXPECT().GetByEmailAndPassword(context.Background(), "success@example.com", mock.AnythingOfType("string")).Return(&entity.User{
					ID:       uuid.New(),
					Email:    "success@example.com",
					Password: "password",
					Name:     "test",
				}, nil).Times(1)
			},
			assertFunc: func(resp *LoginRespDTO, err error) {
				s.Assert().NoError(err)
				s.Assert().Equal("success@example.com", resp.Email)
				s.Assert().Equal("test", resp.Name)
				s.Assert().NotEmpty(resp.Token)
				s.Assert().NotEmpty(resp.ExpiresIn)
			},
		},
	}

	for _, tt := range tests {

		s.Run(tt.name, func() {

			tt.mockFunc(s.mockUserRepo)
			resp, err := s.userService.LoginUser(tt.args.ctx, tt.args.req)
			tt.assertFunc(resp, err)
		})
	}
}
