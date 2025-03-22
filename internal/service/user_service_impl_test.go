package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	pb "github.com/itmrchow/todolist-proto/protobuf/user"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/itmrchow/todolist-users/internal/entity"
	"github.com/itmrchow/todolist-users/internal/repository"
)

func TestRegisterTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}

type RegisterTestSuite struct {
	suite.Suite
	userService  pb.UserServiceServer
	mockUserRepo *repository.MockUsersRepository
	input        struct {
		ctx context.Context
		req *pb.RegisterRequest
	}
}

func (s *RegisterTestSuite) SetupTest() {
	userRepo := repository.NewMockUsersRepository(s.T())
	s.mockUserRepo = userRepo
	s.userService = NewUserService(userRepo)
}

func (s *RegisterTestSuite) Test_Register_ExistsByEmail_DbError() {
	// input
	s.input.ctx = context.Background()
	s.input.req = &pb.RegisterRequest{
		Email:    "db_error@example.com",
		Password: "password",
		Name:     "test",
	}

	// mock
	s.mockUserRepo.EXPECT().ExistsByEmail(context.Background(), mock.Anything).Return(false, errors.New("db error"))

	// execute
	resp, err := s.userService.Register(s.input.ctx, s.input.req)

	// assert
	s.Assert().Nil(resp)
	rpcErr, ok := status.FromError(err)
	s.Assert().True(ok)
	s.Assert().Equal(codes.Internal, rpcErr.Code())

}

func (s *RegisterTestSuite) Test_Register_ExistsByEmail_EmailAlreadyExists() {
	// input
	s.input.ctx = context.Background()
	s.input.req = &pb.RegisterRequest{
		Email:    "exist@example.com",
		Password: "password",
		Name:     "test",
	}

	// mock
	s.mockUserRepo.EXPECT().ExistsByEmail(context.Background(), mock.Anything).Return(true, nil)

	// execute
	resp, err := s.userService.Register(s.input.ctx, s.input.req)

	// assert
	s.Assert().Nil(resp)
	rpcErr, ok := status.FromError(err)
	s.Assert().True(ok)
	s.Assert().Equal(codes.AlreadyExists, rpcErr.Code())
	s.Assert().Equal("email already exists", rpcErr.Message())

}

func (s *RegisterTestSuite) Test_Register_Create_DbError() {
	// input
	s.input.ctx = context.Background()
	s.input.req = &pb.RegisterRequest{
		Email:    "test@example.com",
		Password: "password",
		Name:     "test",
	}

	// mock
	s.mockUserRepo.EXPECT().ExistsByEmail(context.Background(), s.input.req.Email).Return(false, nil)
	s.mockUserRepo.EXPECT().Create(context.Background(), mock.MatchedBy(func(user *entity.User) bool {
		return user.Email == s.input.req.Email &&
			user.Password != s.input.req.Password &&
			user.Name == s.input.req.Name &&
			user.ID != uuid.Nil
	})).Return(errors.New("db error"))

	// execute
	resp, err := s.userService.Register(s.input.ctx, s.input.req)

	// assert
	s.Assert().Nil(resp)
	rpcErr, ok := status.FromError(err)
	s.Assert().True(ok)
	s.Assert().Equal(codes.Internal, rpcErr.Code())
}

func (s *RegisterTestSuite) Test_Register_Success() {

	// input
	s.input.ctx = context.Background()
	s.input.req = &pb.RegisterRequest{
		Email:    "test@example.com",
		Password: "password",
		Name:     "test",
	}

	// mock
	s.mockUserRepo.EXPECT().ExistsByEmail(context.Background(), "test@example.com").Return(false, nil)
	s.mockUserRepo.EXPECT().Create(context.Background(), mock.MatchedBy(func(user *entity.User) bool {
		return user.Email == s.input.req.Email &&
			user.Password != s.input.req.Password &&
			user.Name == s.input.req.Name &&
			user.ID != uuid.Nil
	})).Return(nil)

	// execute
	resp, err := s.userService.Register(s.input.ctx, s.input.req)

	// assert
	s.Assert().Nil(resp)
	s.Assert().Nil(err)
}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}

type LoginTestSuite struct {
	suite.Suite
	userService  pb.UserServiceServer
	mockUserRepo *repository.MockUsersRepository
	input        struct {
		ctx context.Context
		req *pb.LoginRequest
	}
}

func (s *LoginTestSuite) SetupTest() {
	userRepo := repository.NewMockUsersRepository(s.T())
	s.mockUserRepo = userRepo
	s.userService = NewUserService(userRepo)
}

func (s *LoginTestSuite) Test_Login_GetByEmailAndPassword_DBError() {
	// input
	s.input.ctx = context.Background()
	s.input.req = &pb.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	// mock
	s.mockUserRepo.EXPECT().GetByEmailAndPassword(context.Background(), s.input.req.Email, mock.MatchedBy(func(password string) bool {
		return password != s.input.req.Password
	})).Return(nil, errors.New("db error"))

	// execute
	resp, err := s.userService.Login(s.input.ctx, s.input.req)

	// assert
	s.Assert().Nil(resp)
	rpcErr, ok := status.FromError(err)
	s.Assert().True(ok)
	s.Assert().Equal(codes.Internal, rpcErr.Code())
}

func (s *LoginTestSuite) Test_Login_UserNotFound() {
	// input
	s.input.ctx = context.Background()
	s.input.req = &pb.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	// mock
	s.mockUserRepo.EXPECT().GetByEmailAndPassword(context.Background(), s.input.req.Email, mock.MatchedBy(func(password string) bool {
		return password != s.input.req.Password
	})).Return(nil, gorm.ErrRecordNotFound)

	// execute
	resp, err := s.userService.Login(s.input.ctx, s.input.req)

	// assert
	s.Assert().Nil(resp)
	rpcErr, ok := status.FromError(err)
	s.Assert().True(ok)
	s.Assert().Equal(codes.Unauthenticated, rpcErr.Code())
}

func (s *LoginTestSuite) Test_Login_Success() {
	// input
	s.input.ctx = context.Background()
	s.input.req = &pb.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	userID := uuid.New()
	user := &entity.User{
		ID:       userID,
		Email:    s.input.req.Email,
		Name:     "test",
		Password: "hashed_password",
	}

	// mock
	s.mockUserRepo.EXPECT().GetByEmailAndPassword(context.Background(), s.input.req.Email, mock.MatchedBy(func(password string) bool {
		return password != s.input.req.Password
	})).Return(user, nil)

	// execute
	resp, err := s.userService.Login(s.input.ctx, s.input.req)

	// assert
	s.Assert().NotNil(resp)
	s.Assert().Nil(err)
	s.Assert().Equal(userID.String(), resp.Id)
	s.Assert().Equal(user.Name, resp.Name)
	s.Assert().Equal(user.Email, resp.Email)
	s.Assert().NotEmpty(resp.Token)
	s.Assert().NotNil(resp.ExpiresIn)
}

// func (s *UserServiceImplTestSuite) Test_userServiceImpl_Login() {
// 	type args struct {
// 		ctx context.Context
// 		req *pb.LoginRequest
// 	}
// 	tests := []struct {
// 		name       string
// 		args       args
// 		mockFunc   func(m *repository.MockUsersRepository)
// 		assertFunc func(resp *pb.LoginResponse, err error)
// 	}{
// 		{
// 			name: "db error",
// 			args: args{
// 				ctx: context.Background(),
// 				req: &pb.LoginRequest{
// 					Email:    "db_error@example.com",
// 					Password: "password",
// 				},
// 			},
// 			mockFunc: func(m *repository.MockUsersRepository) {
// 				m.EXPECT().GetByEmailAndPassword(context.Background(), "db_error@example.com", mock.AnythingOfType("string")).Return(nil, errors.New("db error")).Times(1)
// 			},
// 			assertFunc: func(resp *pb.LoginResponse, err error) {
// 				s.Assert().ErrorIs(err, &mErr.Err500InternalServer)
// 			},
// 		},
// 		{
// 			name: "user not found",
// 			args: args{
// 				ctx: context.Background(),
// 				req: &pb.LoginRequest{
// 					Email:    "not_found@example.com",
// 					Password: "password",
// 				},
// 			},
// 			mockFunc: func(m *repository.MockUsersRepository) {
// 				m.EXPECT().GetByEmailAndPassword(context.Background(), "not_found@example.com", mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound).Times(1)
// 			},
// 			assertFunc: func(resp *pb.LoginResponse, err error) {
// 				s.Assert().ErrorIs(err, &mErr.Err400InvalidLoginInfo)
// 			},
// 		},
// 		{
// 			name: "success",
// 			args: args{
// 				ctx: context.Background(),
// 				req: &pb.LoginRequest{
// 					Email:    "success@example.com",
// 					Password: "password",
// 				},
// 			},
// 			mockFunc: func(m *repository.MockUsersRepository) {
// 				m.EXPECT().GetByEmailAndPassword(context.Background(), "success@example.com", mock.AnythingOfType("string")).Return(&entity.User{
// 					ID:       uuid.New(),
// 					Email:    "success@example.com",
// 					Password: "password",
// 					Name:     "test",
// 				}, nil).Times(1)
// 			},
// 			assertFunc: func(resp *pb.LoginResponse, err error) {
// 				s.Assert().NoError(err)
// 				s.Assert().Equal("success@example.com", resp.Email)
// 				s.Assert().Equal("test", resp.Name)
// 				s.Assert().NotEmpty(resp.Token)
// 				s.Assert().NotEmpty(resp.ExpiresIn)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {

// 		s.Run(tt.name, func() {

// 			tt.mockFunc(s.mockUserRepo)
// 			resp, err := s.userService.Login(tt.args.ctx, tt.args.req)
// 			tt.assertFunc(resp, err)
// 		})
// 	}
// }
