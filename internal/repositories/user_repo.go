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
	Create(ctx context.Context, email, password, name string) (*models.User, error)
}
type userRepo struct {
	db *pgxpool.Pool
}

func NewUser(db *pgxpool.Pool) UserRepo {
	return &userRepo{db: db}
}

func (ur *userRepo) Create(ctx context.Context, email, password, name string) (*models.User, error) {
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
			return &user, ErrDuplicateEmail
		}
		return &user, newRepoErro(err)
	}
	return &user, nil
}
