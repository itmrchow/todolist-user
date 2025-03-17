package service

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/itmrchow/todolist-users/internal/entity"
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
				s.Assert().Equal("internal server error", err.Error())
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
