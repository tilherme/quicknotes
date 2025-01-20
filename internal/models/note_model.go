package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Note struct {
	Id         pgtype.Numeric
	Title      pgtype.Text
	Content    pgtype.Text
	Color      pgtype.Text
	Created_at pgtype.Date
	Updated_at pgtype.Date
}
