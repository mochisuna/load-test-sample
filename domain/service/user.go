package service

import (
	"context"

	"github.com/mochisuna/load-test-sample/domain"
)

type UserService interface {
	Refer(domain.UserID) (*domain.User, error)
	Register(context.Context, *domain.User) error
}
