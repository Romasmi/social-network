package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Provider interface {
	User() *UserRepository
	Profile() *ProfileRepository
}

type UnitOfWork interface {
	Provider
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type UnitOfWorkImpl struct {
	db *pgxpool.Pool
}

func CreateUnitOfWork(db *pgxpool.Pool) UnitOfWork {
	return &UnitOfWorkImpl{db: db}
}

func (u *UnitOfWorkImpl) User() *UserRepository {
	return CreateUserRepository(u.db)
}

func (u *UnitOfWorkImpl) Profile() *ProfileRepository {
	return CreateProfileRepository(u.db)
}

func (u *UnitOfWorkImpl) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	transaction, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			err := transaction.Rollback(ctx)
			if err != nil {
				fmt.Printf("error while transaction rollback: %v", err)
			}
			panic(p)
		}
	}()

	if err := fn(ctx); err != nil {
		localErr := transaction.Rollback(ctx)
		if localErr != nil {
			fmt.Printf("error while transaction rollback: %v", localErr)
		}
		return err
	}

	return transaction.Commit(ctx)
}
