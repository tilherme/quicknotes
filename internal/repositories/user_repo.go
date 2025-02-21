package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gitgub.com/tilherme/quicknotes/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrDuplicateEmail = newRepoErro(errors.New("Duplicate email"))
var ErrInvalidToken = newRepoErro(errors.New("invalid token or user already confirmed"))

type UserRepo interface {
	Create(ctx context.Context, email, password, name, hashKey string) (*models.User, string, error)
	ConfirmUserByToken(ctx context.Context, token string) error
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

func (ur *userRepo) ConfirmUserByToken(ctx context.Context, token string) error {
	query := `SELECT u.id u_id, t.id t_id FROM users u INNER JOIN users_confirmation_tokens t
			  ON u.id = t.user_id
			  WHERE u.active = false
			  AND t.confirmed = false
			  AND t.token = $1`
	row := ur.db.QueryRow(ctx, query, token)
	var userId, tokenId pgtype.Numeric
	err := row.Scan(&userId, &tokenId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return ErrInvalidToken
		}
		return newRepoErro(err)
	}
	fmt.Println(userId, tokenId, "-------------------")
	queryUpadateUser := "UPDATE users SET active = true, updated_at = now() WHERE id = $1"
	_, err = ur.db.Exec(ctx, queryUpadateUser, userId)
	if err != nil {
		return newRepoErro(err)
	}
	queryUpadateToken := "UPDATE users_confirmation_tokens SET confirmed = true, updated_at = now() WHERE id = $1"
	_, err = ur.db.Exec(ctx, queryUpadateToken, tokenId)
	if err != nil {
		return newRepoErro(err)
	}
	return nil
}
