package userservice

import (
	"fmt"
	"gameapp/entity"
	phonenumber "gameapp/pkg"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
}

type Service struct {
	repo Repository
}

type RegisterRequest struct {
	Name        string
	PhoneNumber string
}

type RegisterResponse struct {
	User entity.User
}

func New(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	// TODO: We should verify phone number by verification code

	// validate phone number

	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("phone number is not valid")
	}

	// check uniqueness of phone number

	if isUnique, error := s.repo.IsPhoneNumberUnique(req.PhoneNumber); error != nil || !isUnique {
		if error != nil {
			return RegisterResponse{}, fmt.Errorf("unexpected error %w", error)
		}

		return RegisterResponse{}, fmt.Errorf("phone number is not unique")
	}

	// validate name

	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("name length should be greater than 3")
	}

	user, error := s.repo.Register(entity.User{PhoneNumber: req.PhoneNumber, Name: req.Name})

	if error != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error: %w", error)
	}

	return RegisterResponse{User: user}, nil
}
