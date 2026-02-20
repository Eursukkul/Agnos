package service

import (
	"errors"
	"strings"
	"time"

	"agnos/internal/model"
	"agnos/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type StaffService interface {
	Create(username, password, hospital string) (model.Staff, error)
	Login(username, password, hospital string) (string, error)
}

type staffService struct {
	repo      repository.StaffRepository
	jwtSecret []byte
	tokenTTL  time.Duration
}

func NewStaffService(repo repository.StaffRepository, jwtSecret string, tokenTTL time.Duration) StaffService {
	return &staffService{repo: repo, jwtSecret: []byte(jwtSecret), tokenTTL: tokenTTL}
}

func (s *staffService) Create(username, password, hospital string) (model.Staff, error) {
	username = strings.TrimSpace(username)
	hospital = strings.TrimSpace(hospital)
	if username == "" || strings.TrimSpace(password) == "" || hospital == "" {
		return model.Staff{}, errors.New("username, password, hospital are required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.Staff{}, err
	}
	return s.repo.Create(username, string(hash), hospital)
}

func (s *staffService) Login(username, password, hospital string) (string, error) {
	user, err := s.repo.FindByUsernameAndHospital(strings.TrimSpace(username), strings.TrimSpace(hospital))
	if err != nil {
		return "", ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", ErrInvalidCredentials
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"staff_id": user.ID,
		"hospital": user.Hospital,
		"iat":      now.Unix(),
		"exp":      now.Add(s.tokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
