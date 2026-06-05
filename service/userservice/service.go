package userservice

import (
	"fmt"
	"gameapp/entity"
	"gameapp/pkg/phonenumber"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
	GetByPhoneNumber(phoneNumber string) (entity.User, bool, error)
	GetById(id uint) (entity.User, error)
}

type AuthService interface {
	CreateAccessToken(userId uint) (string, error)
}

type Service struct {
	repo        UserRepository
	authService AuthService
}

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	user entity.User
}

func New(repo UserRepository, authService AuthService) Service {
	return Service{repo, authService}
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	// TODO: We should verify phone number by verification code

	if req.Name == "" {
		return RegisterResponse{}, fmt.Errorf("name is required")
	}

	if req.PhoneNumber == "" {
		return RegisterResponse{}, fmt.Errorf("phone_number is required")
	}

	if req.Password == "" {
		return RegisterResponse{}, fmt.Errorf("password is required")
	}

	// validate phone number
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("phone number is not valid")
	}

	// check uniqueness of phone number
	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			return RegisterResponse{}, fmt.Errorf("unexpected error %w", err)
		}

		return RegisterResponse{}, fmt.Errorf("phone number is not unique")
	}

	// validate name
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("name length should be greater than 3")
	}

	hashedPassword, encryptionErr := encryptPassword(req.Password)

	if encryptionErr != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error %w", encryptionErr)
	}

	user, registerErr := s.repo.Register(entity.User{PhoneNumber: req.PhoneNumber, Name: req.Name, Password: string(hashedPassword)})

	if registerErr != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error: %w", registerErr)
	}

	return RegisterResponse{user}, nil
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	if req.PhoneNumber == "" {
		return LoginResponse{}, fmt.Errorf("phone_number is required")
	}

	if req.Password == "" {
		return LoginResponse{}, fmt.Errorf("password is required")
	}

	user, notFound, err := s.repo.GetByPhoneNumber(req.PhoneNumber)
	if notFound {
		return LoginResponse{}, fmt.Errorf("invalid phone number or password")
	}

	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return LoginResponse{}, fmt.Errorf("invalid phone number or password")
	}

	accessToken, err := s.authService.CreateAccessToken(user.Id)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error %w", err)
	}

	return LoginResponse{AccessToken: accessToken}, nil
}

type GetProfileRequest struct {
	UserId uint
}

type GetProfileResponse struct {
	Name string `json:"name"`
}

func (s Service) GetProfile(req GetProfileRequest) (GetProfileResponse, error) {
	user, err := s.repo.GetById(req.UserId)

	if err != nil {
		return GetProfileResponse{}, fmt.Errorf("unexpected error %w", err)
	}

	return GetProfileResponse{
		Name: user.Name,
	}, nil
}

func encryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
