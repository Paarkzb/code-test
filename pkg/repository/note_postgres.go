package repository

import (
	"codetest/internal/model"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NotePostgres struct {
	db *pgxpool.Pool
}

func NewNotePostgres(db *pgxpool.Pool) *NotePostgres {
	return &NotePostgres{db: db}
}

func (r *NotePostgres) Create(userId int, note model.Note) (int, error) {
	tx, err := r.db.Begin(context.Background())
	var noteId int

	if err != nil {
		return 0, err
	}

	query := "INSERT INTO public.note (rf_user_id, body) VALUES($1, $2) RETURNING id"
	err = tx.QueryRow(context.Background(), query, userId, note.Body).Scan(&noteId)
	if err != nil {
		_ = tx.Rollback(context.Background())
		return 0, err
	}

	return noteId, tx.Commit(context.Background())
}

func (r *NotePostgres) GetAll(userId int) ([]model.NoteResponse, error) {
	var notes []model.NoteResponse

	query := `
		SELECT 
			note.id as note_id, 
			note.body as note_description,  
			u.id as user_id,
			u.name as user_name,
			u.username as user_username
		FROM public.note 
		LEFT JOIN public.user as u ON u.id = note.rf_user_id
		WHERE note.deleted=false and note.rf_user_id=$1
		ORDER BY note.created_at DESC`

	rows, err := r.db.Query(context.Background(), query, userId)
	for rows.Next() {
		var note model.NoteResponse
		err = rows.Scan(&note.Id, &note.Body, &note.User.Id, &note.User.Name, &note.User.Username)
		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	return notes, err
}

func (r *NotePostgres) GetById(userId, noteId int) (model.NoteResponse, error) {
	var note model.NoteResponse

	query := `
		SELECT 
			note.id, note.body, u.id, u.name, u.username
		FROM public.note
		LEFT JOIN public.user as u on u.id=note.rf_user_id 
		WHERE note.id=$1 and note.rf_user_id=$2 and note.deleted=false`

	err := r.db.QueryRow(context.Background(), query, noteId, userId).Scan(&note.Id, &note.Body, &note.User.Id, &note.User.Name, &note.User.Username)
	if err != nil {
		return note, err
	}

	return note, err
}

func (r *NotePostgres) Delete(userId, noteId int) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}

	query := "UPDATE public.note SET deleted=true WHERE id=$1 and rf_user_id=$2"
	_, err = tx.Exec(context.Background(), query, noteId, userId)
	if err != nil {
		_ = tx.Rollback(context.Background())
		return err
	}

	return tx.Commit(context.Background())
}

func (r *NotePostgres) Update(userId, noteId int, input model.UpdateNoteInput) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}

	noteSetValues := make([]string, 0)
	noteArgs := make([]interface{}, 0)
	noteArgId := 1

	if input.Body != nil {
		noteSetValues = append(noteSetValues, fmt.Sprintf("body=$%d", noteArgId))
		noteArgs = append(noteArgs, *input.Body)
		noteArgId++
	}

	noteSetValues = append(noteSetValues, fmt.Sprintf("updated_at=$%d", noteArgId))
	noteArgs = append(noteArgs, time.Now())
	noteArgId++

	setQuizQuery := strings.Join(noteSetValues, ", ")

	quizQuery := fmt.Sprintf("UPDATE public.note SET %s WHERE id=$%d and rf_user_id=$%d", setQuizQuery, noteArgId, noteArgId+1)
	noteArgs = append(noteArgs, noteId)

	// logrus.Printf("Update note: %s", quizQuery)

	_, err = tx.Exec(context.Background(), quizQuery, noteArgs...)
	if err != nil {
		_ = tx.Rollback(context.Background())
		return err
	}

	return tx.Commit(context.Background())
}
