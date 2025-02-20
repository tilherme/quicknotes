package repositories

import (
	"context"
	"errors"
	"strings"

	"gitgub.com/tilherme/quicknotes/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrDuplicateEmail = newRepoErro(errors.New("Duplicate email"))

type UserRepo interface {
	Create(ctx context.Context, email, password, name, hashKey string) (*models.User, string, error)
}
type userRepo struct {
	db *pgxpool.Pool
}

func NewUser(db *pgxpool.Pool) UserRepo {
	return &userRepo{db: db}
}

func (ur *userRepo) Create(ctx context.Context, email, password, name, hashKey string) (*models.User, string, error) {
	var user models.User
	user.Email = pgtype.Text{String: email, Valid: true}
	user.Password = pgtype.Text{String: password, Valid: true}
	user.Name = pgtype.Text{String: name, Valid: true}
	query := `INSERT INTO users(email, password, name)
	VALUES($1, $2, $3)
	RETURNING id, created_at`
	row := ur.db.QueryRow(ctx, query, email, password, name)
	if err := row.Scan(&user.Id, &user.Created_at); err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			return &user, "", ErrDuplicateEmail
		}
		return &user, "", newRepoErro(err)
	}
	userToken, err := ur.createConfirmationToken(ctx, &user, hashKey)
	if err != nil {
		return nil, "", err
	}
	return &user, userToken.Token.String, nil
}

func (ur *userRepo) createConfirmationToken(ctx context.Context, user *models.User, token string) (*models.UserConfirmationToken, error) {
	var userToken models.UserConfirmationToken
	userToken.Token = pgtype.Text{String: token, Valid: true}
	userToken.UserId = user.Id
	query := `INSERT INTO users_confirmation_tokens (user_id, token)
	VALUES($1, $2)
	RETURNING id, created_at
	`
	row := ur.db.QueryRow(ctx, query, userToken.UserId, userToken.Token)
	if err := row.Scan(&userToken.Id, &userToken.Created_at); err != nil {
		return nil, err
	}
	return &userToken, nil
}
