package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hacKRD0/trikona_go/pkg/logger"
	"go.uber.org/zap"
)

// JWTService defines the interface for JWT operations
type JWTService interface {
	GenerateToken(userID string, email string, role string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	ExtractClaims(token *jwt.Token) (*Claims, error)
}

// Claims represents the JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secret string
}

// NewJWTService creates a new JWT service
func NewJWTService(secret string) JWTService {
	return &jwtService{
		secret: secret,
	}
}

// GenerateToken creates a new JWT token
func (s *jwtService) GenerateToken(userID string, email string, role string) (string, error) {
	logger.Info("Generating JWT token",
		zap.String("user_id", userID),
		zap.String("email", email),
		zap.String("role", role),
	)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		logger.Error("Failed to sign JWT token", err,
			zap.String("user_id", userID),
			zap.String("email", email),
		)
		return "", err
	}

	logger.Info("JWT token generated successfully",
		zap.String("user_id", userID),
		zap.String("email", email),
	)
	return tokenString, nil
}

// ValidateToken validates a JWT token
func (s *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	logger.Info("Validating JWT token")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Error("Unexpected signing method", nil,
				zap.String("method", token.Method.Alg()),
			)
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		logger.Error("Failed to validate JWT token", err)
		return nil, err
	}
	if !token.Valid {
		logger.Error("Token is not valid", nil)
		return nil, errors.New("invalid token")
	}

	logger.Info("JWT token validated successfully")
	return token, nil
}


// ExtractClaims extracts claims from a JWT token
func (s *jwtService) ExtractClaims(token *jwt.Token) (*Claims, error) {
	logger.Info("Extracting claims from JWT token")

	// Since we used ParseWithClaims, token.Claims should already be of type *Claims.
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		logger.Info("Claims extracted successfully",
			zap.String("user_id", claims.UserID),
			zap.String("email", claims.Email),
			zap.String("role", claims.Role),
		)
		return claims, nil
	}

	logger.Error("Invalid token claims", nil)
	return nil, errors.New("invalid token claims")
}

