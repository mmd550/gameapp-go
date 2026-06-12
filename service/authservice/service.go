package authservice

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"gameapp/config"
	"gameapp/pkg/errormessage"
	"gameapp/pkg/richerror"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	ACCESS_TOKEN_SUBJECT  = "access_token"
	REFRESH_TOKEN_SUBJECT = "refresh_token"
)

type Service struct {
	config config.AuthConfig
}

type CustomClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

func New(config config.AuthConfig) *Service {
	return &Service{config: config}
}

func (s Service) CreateAccessToken(userId uint) (string, error) {
	return s.createToken(userId, ACCESS_TOKEN_SUBJECT)
}

func (s Service) CreateRefreshToken(userId uint) (string, error) {
	return s.createToken(userId, REFRESH_TOKEN_SUBJECT)
}

func (s Service) AuthenticateUser(r *http.Request) (uint, error) {
	op := "AuthService.AuthenticateUser"
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, richerror.New(op).WithKind(richerror.KindForbidden).WithMessage(errormessage.Forbidden).WithMeta(map[string]any{
			"details": "no Authorization header in request headers",
		})
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return 0, richerror.New(op).WithKind(richerror.KindForbidden).WithMessage(errormessage.Forbidden).WithMeta(map[string]any{
			"details": fmt.Sprintf("invalid authorization header format value: %s", authHeader),
		})
	}

	claims, parseErr := s.parseToken(parts[1], ACCESS_TOKEN_SUBJECT)

	if parseErr != nil {
		return 0, parseErr
	}

	userId, conversionErr := strconv.ParseUint(claims.UserID, 10, 32)

	if conversionErr != nil {
		return 0, richerror.New(op).WithKind(richerror.KindForbidden).WithMessage(errormessage.Forbidden).WithMeta(map[string]any{
			"details": fmt.Sprintf("unexpected user id %s", claims.UserID),
		})
	}

	return uint(userId), nil
}

func (s Service) createToken(userId uint, subject string) (string, error) {
	expirationTime := time.Now().Add(
		s.config.AccessTokenExpirationMinutes * time.Minute,
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
			Subject:   subject,
		},
		UserID: fmt.Sprintf("%d", userId),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s Service) parseToken(tokenString string, subject string) (*CustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || claims.Subject != subject {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func generateTokenID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
