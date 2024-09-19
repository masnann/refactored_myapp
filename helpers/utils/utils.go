package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"myapp/config"
	"myapp/repository"
	"time"

	"github.com/golang-jwt/jwt"

	"golang.org/x/crypto/argon2"
)

type UtilsService struct {
	repo repository.Repository
}

func NewUtilsService(repo repository.Repository) UtilsService {
	return UtilsService{
		repo: repo,
	}
}

const (
	saltSize    = 16
	keySize     = 32
	timeCost    = 1
	memory      = 64 * 1024
	parallelism = 2
)

func (u UtilsService) GenerateHash(input string) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(input), salt, timeCost, memory, parallelism, keySize)
	saltAndHash := append(salt, hash...)
	encodedSaltAndHash := base64.RawStdEncoding.EncodeToString(saltAndHash)

	return encodedSaltAndHash, nil
}

func (u UtilsService) CompareHash(hash, input string) (bool, error) {
	decodedSaltAndHash, err := base64.RawStdEncoding.DecodeString(hash)
	if err != nil {
		return false, err
	}

	if len(decodedSaltAndHash) < saltSize {
		return false, errors.New("invalid hash format")
	}

	salt := decodedSaltAndHash[:saltSize]
	existingHash := decodedSaltAndHash[saltSize:]

	computedHash := argon2.IDKey([]byte(input), salt, timeCost, memory, parallelism, keySize)

	if subtle.ConstantTimeCompare(existingHash, computedHash) == 1 {
		return true, nil
	}

	return false, errors.New("input mismatch")
}

func (u UtilsService) GenerateJWT(userID int64, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"email":  email,
		"role":   role,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (u UtilsService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (u UtilsService) GenerateRefreshToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (u UtilsService) ValidateRefreshToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["userID"].(float64)
		if !ok {
			return 0, errors.New("invalid user_id type")
		}
		return int64(userID), nil
	} else {
		return 0, errors.New("invalid refresh token")
	}
}

func (u UtilsService) GenerateOTP(length int) (string, error) {
	otp := make([]byte, length)
	_, err := rand.Read(otp)
	if err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		otp[i] = (otp[i] % 10) + '0'
	}

	return string(otp), nil
}

func (u UtilsService) CompareOTP(otpHash, otp string) (bool, error) {
	decodedSaltAndHash, err := base64.RawStdEncoding.DecodeString(otpHash)
	if err != nil {
		return false, err
	}

	if len(decodedSaltAndHash) < saltSize {
		return false, errors.New("invalid hash format")
	}

	salt := decodedSaltAndHash[:saltSize]
	existingHash := decodedSaltAndHash[saltSize:]

	computedHash := argon2.IDKey([]byte(otp), salt, timeCost, memory, parallelism, keySize)

	if subtle.ConstantTimeCompare(existingHash, computedHash) == 1 {
		return true, nil
	}

	return false, errors.New("otp mismatch")
}
