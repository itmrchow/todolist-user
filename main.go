package main

import (
	"fmt"
	"net"

	"github.com/itmrchow/todolist-proto/protobuf/user"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"

	"github.com/itmrchow/todolist-user/internal/infra"
	"github.com/itmrchow/todolist-user/internal/repository"
	"github.com/itmrchow/todolist-user/internal/service"
)

func main() {

	// init config
	initConfig()

	// db conn
	mysqlConn := initMysqlDb()

	// repo
	repo := repository.NewUsersRepository(mysqlConn)

	// grpc
	log.Fatal().Err(RunGrpcHandler(repo)).Msg("failed to listen")
}

func initConfig() {
	err := infra.InitConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init config")
	}

	log.Info().Msg("config loaded")
}

func initMysqlDb() *gorm.DB {
	db, err := infra.InitMysqlDb()

	if err != nil {
		log.Fatal().Err(err).Msg("failed to init mysql db")
	}

	log.Info().Msg("mysql db connected")

	return db
}

func RunGrpcHandler(userRepo repository.UsersRepository) (err error) {

	var (
		grpcPort = viper.GetString("server_port")
	)

	log.Info().Msg("grpc server listen in port:" + grpcPort)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
		return
	}
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
		// TODO: add interceptor
		),
	}

	// user service impl
	userService := service.NewUserService(userRepo)

	s := grpc.NewServer(opts...)
	user.RegisterUserServiceServer(s, userService)

	reflection.Register(s)
	err = s.Serve(lis)

	return err

	// user.RegisterUserServiceServer(grpcServer, &service.UserServiceImpl{})

}
