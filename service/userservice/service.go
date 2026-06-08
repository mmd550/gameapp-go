package userservice

import (
	"fmt"
	"gameapp/dto"
	"gameapp/entity"
	"gameapp/pkg/richerror"

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

func New(repo UserRepository, authService AuthService) Service {
	return Service{repo, authService}
}

func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
	// TODO - We should verify phone number by verification code

	// check uniqueness of phone number
	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			return dto.RegisterResponse{}, fmt.Errorf("unexpected error %w", err)
		}

		return dto.RegisterResponse{}, fmt.Errorf("phone number is not unique")
	}

	hashedPassword, encryptionErr := encryptPassword(req.Password)

	if encryptionErr != nil {
		return dto.RegisterResponse{}, fmt.Errorf("unexpected error %w", encryptionErr)
	}

	user, registerErr := s.repo.Register(entity.User{PhoneNumber: req.PhoneNumber, Name: req.Name, Password: string(hashedPassword)})

	if registerErr != nil {
		return dto.RegisterResponse{}, fmt.Errorf("unexpected error: %w", registerErr)
	}

	return dto.RegisterResponse{User: dto.RegisteredUser{
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Id:          user.Id,
	}}, nil
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (s Service) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	op := "userservice.login"
	if req.PhoneNumber == "" {
		return dto.LoginResponse{}, richerror.New(op).WithKind(richerror.KindInvalid).WithMessage("phone number is required")
	}

	if req.Password == "" {
		return dto.LoginResponse{}, richerror.New(op).WithKind(richerror.KindInvalid).WithMessage("password is required")
	}

	user, notFound, err := s.repo.GetByPhoneNumber(req.PhoneNumber)
	if notFound {
		return dto.LoginResponse{}, richerror.New(op).WithKind(richerror.KindForbidden).WithMessage("invalid phone number or password")
	}

	if err != nil {
		return dto.LoginResponse{}, richerror.New(op).WithErr(err).WithMessage("unexpected error").WithKind(richerror.KindUnexpected)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return dto.LoginResponse{}, richerror.New(op).WithKind(richerror.KindForbidden).WithMessage("invalid phone number or password")
	}

	accessToken, err := s.authService.CreateAccessToken(user.Id)
	if err != nil {
		return dto.LoginResponse{}, richerror.New(op).WithErr(err).WithMessage("unexpected error").WithKind(richerror.KindUnexpected)
	}

	return dto.LoginResponse{AccessToken: accessToken}, nil
}

func (s Service) GetProfile(req dto.GetProfileRequest) (dto.GetProfileResponse, error) {
	user, err := s.repo.GetById(req.UserId)

	if err != nil {
		return dto.GetProfileResponse{}, fmt.Errorf("unexpected error %w", err)
	}

	return dto.GetProfileResponse{
		Name: user.Name,
	}, nil
}

func encryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
