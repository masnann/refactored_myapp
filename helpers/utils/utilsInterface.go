package utils

import (
	"github.com/golang-jwt/jwt"
)

type UtilsInterface interface {
	GenerateHash(input string) (string, error)
	CompareHash(hash, input string) (bool, error)
	GenerateJWT(userID int64, email, role string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GenerateRefreshToken(userID int64) (string, error)
	ValidateRefreshToken(tokenString string) (int64, error)
	GenerateOTP(length int) (string, error)
	CompareOTP(otpHash, otp string) (bool, error)
}
