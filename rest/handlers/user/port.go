package user

import "scheduler-app/domain"

type Service interface {
	Signup(username, password string) (int, error)
	Login(username, password string) (*domain.User, error)
}