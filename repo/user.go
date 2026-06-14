package repo

import (
	"scheduler-app/domain"
	"scheduler-app/user"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	user.UserRepo
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) CreateUser(username, passwordHash string) (int, error) {
	var id int
	query := `
	INSERT INTO users(
		username,
		password_hash
	) 
	VALUES($1, $2)
	RETURNING id`

	row := r.db.QueryRow(query, username, passwordHash)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *userRepo) GetUserByUsername(username string) (*domain.User, error) {
	var u domain.User
	query := `SELECT id, username, password_hash FROM users WHERE username = $1`
	err := r.db.QueryRow(query, username).Scan(&u.ID, &u.Username, &u.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &u, nil
}