package service

import (
	"errors"
	"task/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Register(username, password string) (int, error) {
	if username == "" || password == "" {
		return 0, ErrBadCredentials
	}
	passwordHash, err := s.GetPasswordHash(password)
	if err != nil {
		return 0, errors.New("with GetPasswordHash " + err.Error())
	}
	id, err := s.repo.CreateUser(username, passwordHash)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return 0, ErrUserAlreadyExists
		}
	}
	return id, nil
}

func (s *Service) Login(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", ErrBadCredentials
	}
	
	id, passwordHash,err := s.repo.GetUser(username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
   			return "", ErrUserNotFound
  	  }
		return "", errors.New("with repo.GetUser " + err.Error())
	}
	if !s.ComparePasswords(password, passwordHash) {
		return "", ErrWrongPassword
	}
	
	token, err := s.generateToken(id)
	if err != nil {
		return "", errors.New("with generateToken " + err.Error())
	}
	return token, nil
}

func (s *Service) generateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.cfg.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}


func (s *Service) GetPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *Service) ComparePasswords(password, passwordHash string) bool {
	bytePass := []byte(password)
	bytePassHash := []byte(passwordHash)
	
	if err := bcrypt.CompareHashAndPassword(bytePassHash, bytePass); err != nil {
		return false
	}
	return true
}

func (s *Service) ParseToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}
	return []byte(s.cfg.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return 0, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user ID in token claims")
	}
	expTime, ok := claims["exp"].(float64)
 	if !ok || time.Unix(int64(expTime), 0).Before(time.Now()) {
  		return 0, errors.New("token has expired")
 	}
	return int(userID), nil
}