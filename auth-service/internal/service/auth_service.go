package service

import (
	"errors"
	"time"

	"github.com/Bkgediya/feed_system/auth-service/internal/auth"
	"github.com/Bkgediya/feed_system/auth-service/internal/model"
	"github.com/Bkgediya/feed_system/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService interface {
	SignUp(username, email, password string) (*model.User, error)
	Login(email, password string) (string, error)
	GetUser(id int64) (*model.User, error)
}

type authService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewAuthService(repo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{repo: repo, jwtSecret: jwtSecret}
}

func (s *authService) SignUp(username, email, password string) (*model.User, error) {
	existing, _ := s.repo.GetByEmail(email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *authService) Login(email, password string) (string, error) {
	u, err := s.repo.GetByEmail(email)
	if err != nil || u == nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := auth.GenerateJWT(u.ID, u.Email, s.jwtSecret, time.Hour*24)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) GetUser(id int64) (*model.User, error) {
	return s.repo.GetByID(id)
}
