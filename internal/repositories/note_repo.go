package repositories

import (
	"context"
	"math/big"
	"time"

	"gitgub.com/tilherme/quicknotes/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NoteRepo interface {
	List(ctx context.Context) ([]models.Note, error)
	GetById(ctx context.Context, id int) (*models.Note, error)
	Create(ctx context.Context, title, content, color string) (*models.Note, error)
	Update(ctx context.Context, id int, title, content, color string) (*models.Note, error)
	Delete(ctx context.Context, id int) error
}

func NewNote(dbpool *pgxpool.Pool) NoteRepo {
	return &noteRepo{db: dbpool}
}

type noteRepo struct {
	db *pgxpool.Pool
}

func (nr *noteRepo) List(ctx context.Context) ([]models.Note, error) {
	var list []models.Note
	rows, err := nr.db.Query(ctx, "select * from notes")
	if err != nil {
		return list, newRepoErro(err)
	}
	defer rows.Close()

	for rows.Next() {
		var note models.Note
		err = rows.Scan(&note.Id, &note.Title, &note.Content, &note.Color, &note.Created_at, &note.Updated_at)
		if err != nil {
			return list, newRepoErro(err)
		}
		list = append(list, note)
	}
	return list, nil
}

func (nr *noteRepo) GetById(ctx context.Context, id int) (*models.Note, error) {
	var note models.Note

	row := nr.db.QueryRow(ctx,
		`select * from notes where id = $1`, id)
	if err := row.Scan(&note.Id, &note.Title, &note.Content, &note.Color, &note.Created_at, &note.Updated_at); err != nil {
		return &note, newRepoErro(err)

	}
	return &note, nil
}

func (nr *noteRepo) Create(ctx context.Context, title, content, color string) (*models.Note, error) {
	var note models.Note
	note.Title = pgtype.Text{String: title, Valid: true}
	note.Content = pgtype.Text{String: content, Valid: true}
	note.Color = pgtype.Text{String: color, Valid: true}
	query := `INSERT INTO notes(title, content, color)
	VALUES($1, $2, $3)
	RETURNING id, created_at`
	row := nr.db.QueryRow(ctx, query, title, content, color)
	if err := row.Scan(&note.Id, &note.Created_at); err != nil {
		return &note, newRepoErro(err)
	}
	return &note, nil
}

func (nr *noteRepo) Update(ctx context.Context, id int, title, content, color string) (*models.Note, error) {
	var note models.Note
	note.Id = pgtype.Numeric{Int: big.NewInt(int64(id)), Valid: true}
	if len(title) > 0 {
		note.Title = pgtype.Text{String: title, Valid: true}
	}
	if len(content) > 0 {
		note.Content = pgtype.Text{String: content, Valid: true}
	}
	if len(color) > 0 {
		note.Color = pgtype.Text{String: color, Valid: true}
	}
	note.Updated_at = pgtype.Date{Time: time.Now(), Valid: true}
	query := `UPDATE notes set title = $1 ,content = COALESCE($2, content), color = $3, updated_at = $4 where id = $5`

	_, err := nr.db.Exec(ctx, query, title, content, color, note.Updated_at, note.Id)
	if err != nil {
		return &note, newRepoErro(err)
	}
	return &note, nil
}

func (nr *noteRepo) Delete(ctx context.Context, id int) error {
	_, err := nr.db.Exec(ctx, "DELETE FROM notes where id = $1;", id)
	if err != nil {
		return newRepoErro(err)
	}
	return nil
}
