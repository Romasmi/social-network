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
	WithTransaction(ctx context.Context, fn func(ctx context.Context, uof UnitOfWork) error) error
}

type UnitOfWorkImpl struct {
	writer DBQuerier
	reader DBQuerier
}

func CreateUnitOfWork(writer DBQuerier, reader DBQuerier) UnitOfWork {
	return &UnitOfWorkImpl{writer: writer, reader: reader}
}

func (u *UnitOfWorkImpl) User() *UserRepository {
	return CreateUserRepository(u.writer, u.reader)
}

func (u *UnitOfWorkImpl) Profile() *ProfileRepository {
	return CreateProfileRepository(u.writer, u.reader)
}

func (u *UnitOfWorkImpl) WithTransaction(ctx context.Context, fn func(ctx context.Context, txUoW UnitOfWork) error) error {
	pool, ok := u.writer.(*pgxpool.Pool)
	if !ok {
		return fn(ctx, u)
	}

	transaction, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			err := transaction.Rollback(ctx)
			if err != nil {
				// TODO replace with logger
				fmt.Printf("error while transaction rollback: %v", err)
			}
			panic(p)
		}
	}()

	txUoW := &UnitOfWorkImpl{writer: transaction, reader: transaction}
	if err := fn(ctx, txUoW); err != nil {
		localErr := transaction.Rollback(ctx)
		if localErr != nil {
			fmt.Printf("error while transaction rollback: %v", localErr)
		}
		return err
	}

	return transaction.Commit(ctx)
}
