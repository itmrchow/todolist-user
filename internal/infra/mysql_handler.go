package infra

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/itmrchow/todolist-users/internal/entity"
)

func InitMysqlDb() (*gorm.DB, error) {

	db, err := gorm.Open(mysql.Open(getDNS()), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()

	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Minute * 30)

	return db, nil
}

func getDNS() (dns string) {

	// account:password@tcp(host:port)/{db_name}{url_suffix}

	dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s",
		viper.GetString("mysql.db_account"),
		viper.GetString("mysql.db_password"),
		viper.GetString("mysql.db_host"),
		viper.GetString("mysql.db_port"),
		viper.GetString("mysql.db_name"),
		viper.GetString("mysql.url_suffix"))

	return
}
