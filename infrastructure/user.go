package infrastructure

import (
	"context"
	"database/sql"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/mochisuna/load-test-sample/domain"
	"github.com/mochisuna/load-test-sample/domain/repository"
	"github.com/mochisuna/load-test-sample/infrastructure/db"
)

func NewUserRepository(dbmClient *db.Client, dbsClient *db.Client) repository.UserRepository {
	return &userRepository{
		dbm: dbmClient,
		dbs: dbsClient,
	}
}

type userRepository struct {
	dbm *db.Client
	dbs *db.Client
}

func (r *userRepository) WithTransaction(ctx context.Context, txFunc func(*sql.Tx) error) error {
	tx, err := r.dbm.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
}

func (r *userRepository) Create(user *domain.User) error {
	log.Println("called infrastructure Create")
	// Exec Create
	_, err := squirrel.Insert("users").
		Columns("id", "name", "secret_key", "created_at", "updated_at").
		Values(user.ID, user.Name, user.SecretKey, user.CreatedAt, user.UpdatedAt).
		RunWith(r.dbm.DB).
		Exec()
	return err
}

type userColumns struct {
	ID        domain.UserID `db:"id"`
	Name      string        `db:"name"`
	SecretKey string        `db:"secret_key"`
	CreatedAt int           `db:"created_at"`
	UpdatedAt int           `db:"updated_at"`
}

func (r *userRepository) Get(userID domain.UserID) (*domain.User, error) {
	log.Println("called infrastructure Get")
	ret := &domain.User{}
	var user userColumns
	err := squirrel.Select("id", "name", "secret_key", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{
			"users.id": userID,
		}).
		RunWith(r.dbs.DB).
		QueryRow().
		Scan(
			&user.ID,
			&user.Name,
			&user.SecretKey,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
	if err != nil {
		return ret, err
	}
	ret = &domain.User{
		ID:        user.ID,
		Name:      user.Name,
		SecretKey: user.SecretKey,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return ret, err
}
