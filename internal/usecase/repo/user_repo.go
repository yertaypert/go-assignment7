package repo

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/yertaypert/go-assignment7/internal/entity"
	"github.com/yertaypert/go-assignment7/pkg"
)

type UserRepo struct {
	PG *pkg.Postgres
}

func NewUserRepo(pg *pkg.Postgres) *UserRepo {
	return &UserRepo{PG: pg}
}

func (u *UserRepo) RegisterUser(user *entity.User) (*entity.User, error) {
	email := strings.ToLower(strings.TrimSpace(user.Email))
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	if err := u.PG.Conn.QueryRow(checkQuery, email).Scan(&exists); err != nil {
		return nil, fmt.Errorf("check existing user: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user with email %s already exists", user.Email)
	}

	created := *user
	created.ID = uuid.New()
	created.Email = email
	if created.Role == "" {
		created.Role = "user"
	}
	created.Verified = false

	const insertQuery = `
INSERT INTO users (id, username, email, password, role, verified)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, username, email, password, role, verified`

	row := u.PG.Conn.QueryRow(
		insertQuery,
		created.ID,
		created.Username,
		created.Email,
		created.Password,
		created.Role,
		created.Verified,
	)

	var stored entity.User
	var id string
	if err := row.Scan(
		&id,
		&stored.Username,
		&stored.Email,
		&stored.Password,
		&stored.Role,
		&stored.Verified,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user insert returned no rows")
		}
		return nil, fmt.Errorf("insert user: %w", err)
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("parse returned id: %w", err)
	}
	stored.ID = parsedID

	return &stored, nil
}

func (u *UserRepo) LoginUser(user *entity.LoginUserDTO) (*entity.User,
	error) {
	username := strings.TrimSpace(user.Username)
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	const query = `
SELECT id, username, email, password, role, verified
FROM users
WHERE username = $1
LIMIT 1`

	var userFromDB entity.User
	var id string
	err := u.PG.Conn.QueryRow(query, username).Scan(
		&id,
		&userFromDB.Username,
		&userFromDB.Email,
		&userFromDB.Password,
		&userFromDB.Role,
		&userFromDB.Verified,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("username not found")
		}
		return nil, fmt.Errorf("find user by username: %w", err)
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}
	userFromDB.ID = parsedID

	return &userFromDB, nil
}
