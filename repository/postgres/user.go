package postgres

import (
	"context"
	"errors"
	"fmt"

	"gameapp/entity"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string
	PhoneNumber string
	Password    string
}

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	ctx := context.Background()

	_, err := gorm.G[User](r.db.conn).Where("phone_number = ?", phoneNumber).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		return false, err
	}

	return false, nil
}

func (r *UserRepository) Register(u entity.User) (entity.User, error) {
	ctx := context.Background()

	newUser := User{Name: u.Name, PhoneNumber: u.PhoneNumber, Password: u.Password}
	err := gorm.G[User](r.db.conn).Create(ctx, &newUser)
	if err != nil {
		return entity.User{}, err
	}

	return entity.User{Id: newUser.ID, Name: newUser.Name, PhoneNumber: newUser.PhoneNumber}, nil
}

func (r *UserRepository) GetByPhoneNumber(phoneNumber string) (entity.User, bool, error) {
	ctx := context.Background()

	user, err := gorm.G[User](r.db.conn).
		Where("phone_number = ?", phoneNumber).
		First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, true, fmt.Errorf("user not found with phone: %s", phoneNumber)
		}
		return entity.User{}, false, fmt.Errorf("failed to get user by phone: %w", err)
	}

	return entity.User{
		Id:          user.ID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Password:    user.Password,
	}, false, nil
}
