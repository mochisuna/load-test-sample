package repository

import (
	"context"
	"database/sql"

	"github.com/mochisuna/load-test-sample/domain"
)

type UserRepository interface {
	WithTransaction(ctx context.Context, txFunc func(*sql.Tx) error) error
	Get(domain.UserID) (*domain.User, error)
	Create(*domain.User) error
}
