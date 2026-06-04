package postgres

import (
	"context"
	"errors"

	"gameapp/entity"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string
	PhoneNumber string
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

	newUser := User{Name: u.Name, PhoneNumber: u.PhoneNumber}
	err := gorm.G[User](r.db.conn).Create(ctx, &newUser)
	if err != nil {
		return entity.User{}, err
	}

	return entity.User{Id: newUser.ID, Name: newUser.Name, PhoneNumber: newUser.PhoneNumber}, nil
}
