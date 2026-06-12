package postgres

import (
	"context"
	"errors"

	"gameapp/entity"
	"gameapp/pkg/errormessage"
	"gameapp/pkg/richerror"

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

func (r *UserRepository) DoesPhoneNumberExist(phoneNumber string) (bool, error) {
	ctx := context.Background()

	op := "UserRepository.DoesPhoneNumberExist"

	_, err := gorm.G[User](r.db.conn).Where("phone_number = ?", phoneNumber).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		return false, richerror.New(op).WithErr(err).WithMessage(errormessage.SomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	return false, nil
}

func (r *UserRepository) Register(u entity.User) (entity.User, error) {
	ctx := context.Background()

	op := "UserRepository.Register"

	newUser := User{Name: u.Name, PhoneNumber: u.PhoneNumber, Password: u.Password}
	err := gorm.G[User](r.db.conn).Create(ctx, &newUser)
	if err != nil {
		return entity.User{}, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected).WithMessage(errormessage.SomethingWentWrong)
	}

	return entity.User{Id: newUser.ID, Name: newUser.Name, PhoneNumber: newUser.PhoneNumber}, nil
}

func (r *UserRepository) GetByPhoneNumber(phoneNumber string) (entity.User, error) {
	ctx := context.Background()

	op := "UserRepository.GetByPhoneNumber"

	user, err := gorm.G[User](r.db.conn).
		Where("phone_number = ?", phoneNumber).
		First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, richerror.New(op).WithKind(richerror.KindNotFound).WithErr(err).WithMessage(errormessage.NotFound)
		}
		return entity.User{}, richerror.New(op).WithKind(richerror.KindUnexpected).WithErr(err).WithMessage(errormessage.SomethingWentWrong)
	}

	return entity.User{
		Id:          user.ID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Password:    user.Password,
	}, nil
}

func (r *UserRepository) GetById(id uint) (entity.User, error) {
	ctx := context.Background()

	op := "UserRepository.GetById"

	user, err := gorm.G[User](r.db.conn).Where("id = ?", id).First(ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, richerror.New(op).WithErr(err).WithMessage(errormessage.NotFound).WithKind(richerror.KindNotFound)
		}
		return entity.User{}, richerror.New(op).WithErr(err).WithMessage(errormessage.SomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	return entity.User{
		Id:          user.ID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Password:    user.Password,
	}, nil
}
