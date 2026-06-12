package userservice

import (
	"fmt"
	"gameapp/dto"
	"gameapp/entity"
	"gameapp/pkg/errormessage"
	"gameapp/pkg/httpmessage"
	"gameapp/pkg/richerror"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Register(u entity.User) (entity.User, error)
	GetByPhoneNumber(phoneNumber string) (entity.User, error)
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
	op := "UserService.Login"

	user, err := s.repo.GetByPhoneNumber(req.PhoneNumber)

	if err != nil {
		_, status := httpmessage.Error(err)
		if status == http.StatusNotFound {
			return dto.LoginResponse{}, richerror.New(op).WithKind(richerror.KindForbidden).WithMessage(errormessage.InvalidPhoneOrPassword)
		}
		return dto.LoginResponse{}, richerror.New(op).WithErr(err).WithMessage("unexpected error").WithKind(richerror.KindUnexpected)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return dto.LoginResponse{}, richerror.New(op).WithKind(richerror.KindForbidden).WithMessage(errormessage.InvalidPhoneOrPassword)
	}

	accessToken, err := s.authService.CreateAccessToken(user.Id)
	if err != nil {
		return dto.LoginResponse{}, richerror.New(op).WithErr(err).WithMessage("unexpected error").WithKind(richerror.KindUnexpected)
	}

	return dto.LoginResponse{AccessToken: accessToken}, nil
}

func (s Service) GetProfile(req dto.GetProfileRequest) (dto.GetProfileResponse, error) {
	op := "UserService.GetProfile"
	user, err := s.repo.GetById(req.UserId)

	if err != nil {
		return dto.GetProfileResponse{}, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected).WithMessage(errormessage.SomethingWentWrong)
	}

	return dto.GetProfileResponse{
		Name: user.Name,
	}, nil
}

func encryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
