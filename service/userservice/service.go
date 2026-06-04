package userservice

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"gameapp/config"
	"gameapp/entity"
	"gameapp/pkg/phonenumber"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
	GetByPhoneNumber(phoneNumber string) (entity.User, bool, error)
}

type Service struct {
	repo      UserRepository
	jwtConfig config.JWTConfig
}

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	user entity.User
}

func New(repo UserRepository, jwtConfig config.JWTConfig) Service {
	return Service{repo, jwtConfig}
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

	accessToken, err := s.createToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error %w", err)
	}

	return LoginResponse{AccessToken: accessToken}, nil
}

type CustomClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

func (s Service) ParseToken(r *http.Request) (*CustomClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header is required")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return nil, fmt.Errorf("authorization header format must be: Bearer <token>")
	}

	token, err := jwt.ParseWithClaims(parts[1], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtConfig.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (s Service) createToken(user entity.User) (string, error) {
	expirationTime := time.Now().Add(
		time.Duration(s.jwtConfig.ExpirationDurationMinutes) * time.Minute,
	)

	jti, err := generateTokenID()
	if err != nil {
		return "", fmt.Errorf("failed to generate token ID: %w", err)
	}

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gameapp",
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        jti, // Protect against replay attacks
		},
		UserID: fmt.Sprintf("%d", user.Id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func encryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// generateTokenID creates a unique identifier for each token
func generateTokenID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
