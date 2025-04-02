package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/itmrchow/todolist-user/internal/entity"
)

var _ UsersRepository = &database{}

type database struct {
	conn *gorm.DB
}

func NewUsersRepository(conn *gorm.DB) UsersRepository {
	return &database{
		conn: conn,
	}
}

func (d *database) Create(ctx context.Context, user *entity.User) error {
	return d.conn.WithContext(ctx).Create(user).Error
}

func (d *database) Get(ctx context.Context, id uuid.UUID) (user *entity.User, err error) {
	if err := d.conn.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return
}

func (d *database) GetByEmail(ctx context.Context, email string) (user *entity.User, err error) {
	if err := d.conn.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return
}

func (d *database) Update(ctx context.Context, user *entity.User) error {
	return d.conn.WithContext(ctx).Model(&entity.User{}).Where("id = ?", user.ID).Updates(user).Error
}

func (d *database) Delete(ctx context.Context, id uuid.UUID) error {
	return d.conn.WithContext(ctx).Where("id = ?", id).Delete(&entity.User{}).Error
}

func (d *database) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := d.conn.WithContext(ctx).Model(&entity.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *database) GetByEmailAndPassword(ctx context.Context, email, password string) (user *entity.User, err error) {
	if err := d.conn.WithContext(ctx).Where("email = ? AND password = ?", email, password).First(&user).Error; err != nil {
		return nil, err
	}
	return
}
