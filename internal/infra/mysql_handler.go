package infra

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/itmrchow/todolist-users/internal/entity"
)

func InitMysqlDb(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		panic("failed to migrate database")
	}

	sqlDB, err := db.DB()

	if err != nil {
		panic("failed to get sql db")
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Minute * 30)

	return db, nil
}
