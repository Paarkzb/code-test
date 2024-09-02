package service

import (
	"codetest/internal/model"
	"codetest/pkg/repository"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(userId int) (model.UserResponse, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}

type Note interface {
	Create(userId int, note model.Note) (int, error)
	GetAll(userId int) ([]model.NoteResponse, error)
	GetById(userId, noteId int) (model.NoteResponse, error)
	Delete(userId, noteId int) error
	Update(userId, noteId int, input model.UpdateNoteInput) error
}

type Service struct {
	Authorization
	Note
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Note:          NewNoteService(repos.Note),
	}
}
