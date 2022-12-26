package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	snsgo "github.com/rifqoi/sns-go"
	"github.com/rifqoi/sns-go/postgres/dbsqlc"
)

type UserRepository struct {
	q *dbsqlc.Queries
}

func New(db dbsqlc.DBTX) *UserRepository {
	return &UserRepository{
		q: dbsqlc.New(db),
	}
}

func (pr *UserRepository) AddUser(ctx context.Context, u snsgo.User) (*snsgo.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := pr.q.AddUser(ctx, dbsqlc.AddUserParams{
		ID:          u.GetNewID(),
		Username:    u.GetUsername(),
		Name:        u.GetName(),
		Password:    u.Password,
		Email:       u.Email,
		Age:         int32(u.Age),
		Gender:      newGender(u.Gender),
		DateOfBirth: u.DateOfBirth,
	})
	if err != nil {
		err = checkDuplicate(err, u)
		return nil, snsgo.WrapErrorf(err, snsgo.ErrorCodeUnknown, "failed to perform db transaction")
	}

	response := &snsgo.User{
		ID:          result.ID,
		Username:    result.Username,
		Password:    result.Password,
		Name:        result.Name,
		Age:         int(result.Age),
		Email:       result.Email,
		Gender:      convertGender(result.Gender),
		DateOfBirth: result.DateOfBirth,
	}

	return response, nil
}

func (u *UserRepository) GetUsers(ctx context.Context) ([]snsgo.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

	defer cancel()

	users, err := u.q.GetUsers(ctx)
	if err != nil {
		return nil, snsgo.WrapErrorf(err, snsgo.ErrorCodeUnknown, "failed to perform db transaction")
	}

	var respUsers []snsgo.User
	for _, user := range users {
		respUsers = append(respUsers, snsgo.User{
			ID:          user.ID,
			Username:    user.Username,
			Name:        user.Name,
			Gender:      convertGender(user.Gender),
			DateOfBirth: user.DateOfBirth,
			Age:         int(user.Age),
			Email:       user.Email,
		})
	}

	return respUsers, nil
}

func (u *UserRepository) GetUserByID(ctx context.Context, ID uuid.UUID) (*snsgo.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := u.q.GetUserByID(ctx, ID)
	if err != nil {
		return nil, snsgo.WrapErrorf(err, snsgo.ErrorCodeUnknown, "failed to perform db transaction")
	}

	respUser := &snsgo.User{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Age:         int(user.Age),
		Email:       user.Email,
		Gender:      convertGender(user.Gender),
		DateOfBirth: user.DateOfBirth,
	}

	return respUser, nil
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*snsgo.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := u.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, snsgo.NewErrorf(snsgo.ErrorCodeUnknown, "user with email %s not found", email)
	}

	respUser := &snsgo.User{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Age:         int(user.Age),
		Email:       user.Email,
		Gender:      convertGender(user.Gender),
		DateOfBirth: user.DateOfBirth,
	}

	return respUser, nil
}

func checkDuplicate(err error, user snsgo.User) error {
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if strings.Contains(pgErr.Message, "username") {
					err = fmt.Errorf("User with username %s already exists.", user.Username)
				} else if strings.Contains(pgErr.Message, "email") {
					err = fmt.Errorf("User with email %s already exists.", user.Email)
				}
			}
		}
	}

	return err
}

// auth => Login, Logout
// account => Register, Update
// Post => CreatePost, UpdatePost
// comment => NewComment, UpdateComment
// likes => NewLike, RemoveLike
// Share => ShareFromPost,
