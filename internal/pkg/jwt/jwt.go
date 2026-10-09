package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type TokenType string

const (
	TokenTypeAccess       TokenType = "ACCESS"
	TokenTypeVerification TokenType = "VERIFICATION"
	TokenTypeReset        TokenType = "RESET"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secretKey      []byte
	accessExpHours int
}

func NewTokenService(viper *viper.Viper) *TokenService {
	secret := viper.GetString("JWT_SECRET")
	if secret == "" {
		secret = "default_fallback_secret_nexus"
	}
	expHours := viper.GetInt("JWT_EXPIRATION_HOURS")
	if expHours == 0 {
		expHours = 24
	}

	return &TokenService{
		secretKey:      []byte(secret),
		accessExpHours: expHours,
	}
}

func (s *TokenService) GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	return s.generateToken(userID, email, TokenTypeAccess, time.Duration(s.accessExpHours)*time.Hour)
}

func (s *TokenService) GenerateVerificationToken(userID uuid.UUID, email string) (string, error) {
	return s.generateToken(userID, email, TokenTypeVerification, 24*time.Hour)
}

func (s *TokenService) GeneratePasswordResetToken(userID uuid.UUID, email string) (string, error) {
	return s.generateToken(userID, email, TokenTypeReset, 1*time.Hour)
}

func (s *TokenService) generateToken(userID uuid.UUID, email string, tokenType TokenType, duration time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *TokenService) ValidateToken(tokenStr string, expectedType TokenType) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	if claims.Type != expectedType {
		return nil, errors.New("token type mismatch")
	}

	return claims, nil
}
