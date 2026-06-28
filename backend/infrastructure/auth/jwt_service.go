package auth

import (
	"fmt"
	"time"

	"wms/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userID, role string) (string, error)
	ValidateToken(tokenString string) (map[string]interface{}, error)
}

type jwtService struct {
	secret     []byte
	expiry     time.Duration
	refreshExp time.Duration
}

func NewJWTService(cfg config.AuthConfig) JWTService {
	expiry, _ := time.ParseDuration(cfg.JWTExpiry)
	if expiry == 0 {
		expiry = 24 * time.Hour
	}

	refreshExp, _ := time.ParseDuration(cfg.RefreshExpiry)
	if refreshExp == 0 {
		refreshExp = 168 * time.Hour
	}

	return &jwtService{
		secret:     []byte(cfg.JWTSecret),
		expiry:     expiry,
		refreshExp: refreshExp,
	}
}

func (s *jwtService) GenerateToken(userID, role string) (string, error) {
	claims := map[string]interface{}{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(s.expiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	return token.SignedString(s.secret)
}

func (s *jwtService) ValidateToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
