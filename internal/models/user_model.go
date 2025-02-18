package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Id         pgtype.Numeric
	Email      pgtype.Text
	Password   pgtype.Text
	Name       pgtype.Text
	Active     pgtype.Bool
	Created_at pgtype.Date
	Updated_at pgtype.Date
}
