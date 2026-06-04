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

func (d *PostgresDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	db := d.Conn()
	ctx := context.Background()

	_, err := gorm.G[User](db).Where("phone_number = ?", phoneNumber).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		return false, err
	}

	return false, nil
}

func (d *PostgresDB) Register(u entity.User) (entity.User, error) {
	db := d.Conn()
	ctx := context.Background()

	newUser := User{Name: u.Name, PhoneNumber: u.PhoneNumber}
	err := gorm.G[User](db).Create(ctx, &newUser)
	if err != nil {
		return entity.User{}, err
	}

	return entity.User{Id: newUser.ID, Name: newUser.Name, PhoneNumber: newUser.PhoneNumber}, nil
}
