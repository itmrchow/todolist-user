package main

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/itmrchow/todolist-users/internal/infra"
	"github.com/itmrchow/todolist-users/internal/repository"
)

func main() {
	fmt.Println("Start init") // log

	initConfig()

	mysqlConn := initMysqlDb()

	// db repo
	// usersRepo := repository.NewUsersRepository(db)
	repository.NewUsersRepository(mysqlConn)

}

func initConfig() {
	infra.InitConfig()
	println("config loaded") // TODO: log
}

func initMysqlDb() *gorm.DB {
	db, err := infra.InitMysqlDb()

	if err != nil {
		panic(err)
	}

	println("db connected") // TODO: log

	return db
}
