package repository

import (
	"medodstest/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(username, password string) (model.User, error)
	GetUserById(id int) (model.UserResponse, error)
}

type Note interface {
	Create(userId int, note model.Note) (int, error)
	GetAll(userId int) ([]model.NoteResponse, error)
	GetById(userId, noteId int) (model.NoteResponse, error)
	Delete(userId, noteId int) error
	Update(userId, noteId int, input model.UpdateNoteInput) error
}

type Repository struct {
	Authorization
	Note
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Note:          NewNotePostgres(db),
	}
}
