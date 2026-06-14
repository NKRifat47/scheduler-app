package user

import (
	"errors"
	"scheduler-app/domain"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	usrRepo UserRepo
}

func NewService (usrRepo UserRepo) Service{
	return &service{
		usrRepo: usrRepo,
	}
}

func (s *service) Signup(username, password string) (int, error) {
	// 1. Check if username exists
	_, err := s.usrRepo.GetUserByUsername(username)
	if err == nil {
		return 0, errors.New("username is already taken")
	}

	// 2. Hash password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// 3. Create user in repository
	return s.usrRepo.CreateUser(username, string(hashedBytes))
}

func (s *service) Login(username, password string) (*domain.User, error) {
	u, err := s.usrRepo.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	return u, nil
}