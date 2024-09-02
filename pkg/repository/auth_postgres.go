package repository

import (
	"codetest/internal/model"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthPostgres struct {
	db *pgxpool.Pool
}

func NewAuthPostgres(db *pgxpool.Pool) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user model.User) (int, error) {
	var id int
	query := "INSERT INTO public.user (name, username, password) VALUES ($1, $2, $3) RETURNING id;"

	err := r.db.QueryRow(context.Background(), query, user.Name, user.Username, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (model.User, error) {
	var user model.User

	query := `
				SELECT u.id
				FROM public.user as u 
				WHERE u.username = $1 AND u.password = $2
			`

	err := r.db.QueryRow(context.Background(), query, username, password).Scan(&user.Id)

	return user, err
}

func (r *AuthPostgres) GetUserById(id int) (model.UserResponse, error) {
	var user model.UserResponse

	query := `
				SELECT u.id, u.name, u.username
				FROM public.user as u 
				WHERE u.id = $1
	`

	err := r.db.QueryRow(context.Background(), query, id).Scan(&user.Id, &user.Name, &user.Username)

	return user, err
}
