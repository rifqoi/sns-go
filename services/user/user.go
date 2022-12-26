package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	snsgo "github.com/rifqoi/sns-go"
	"github.com/rifqoi/sns-go/postgres"
)

type UserService struct {
	user snsgo.UserRepository
}

type UserConfiguration func(us *UserService) error

func NewUserService(configs ...UserConfiguration) (*UserService, error) {
	us := new(UserService)

	for _, config := range configs {
		err := config(us)
		if err != nil {
			return nil, err
		}
	}

	return us, nil
}

func WithUserRepository(ur snsgo.UserRepository) UserConfiguration {
	return func(us *UserService) error {
		us.user = ur
		return nil
	}
}

func WithPostgresUserRepository(db *pgxpool.Pool) UserConfiguration {
	ur := postgres.New(db)
	return WithUserRepository(ur)
}

func (u *UserService) RegisterUser(ctx context.Context, req snsgo.CreateRequest) error {
	err := req.Validate()
	if err != nil {
		return snsgo.WrapErrorf(err, snsgo.ErrorCodeInvalidArg, "")
	}
	user := req.NewUser()
	_, err = u.user.AddUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetUsers(ctx context.Context) ([]snsgo.User, error) {
	users, err := u.user.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserService) GetUserByID(ctx context.Context, ID uuid.UUID) (*snsgo.User, error) {
	user, err := u.user.GetUserByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) GetUserByEmail(ctx context.Context, email string) (*snsgo.User, error) {
	user, err := u.user.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, err
}
