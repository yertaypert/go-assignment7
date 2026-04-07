package repo

import (
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/yertaypert/go-assignment7/internal/entity"
)

type UserRepo struct {
	mu      sync.Mutex
	byEmail map[string]entity.User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		byEmail: make(map[string]entity.User),
	}
}

func (u *UserRepo) RegisterUser(user *entity.User) (*entity.User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	email := strings.ToLower(strings.TrimSpace(user.Email))
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if _, exists := u.byEmail[email]; exists {
		return nil, fmt.Errorf("user with email %s already exists", user.Email)
	}

	created := *user
	created.ID = uuid.New()
	created.Email = email
	if created.Role == "" {
		created.Role = "user"
	}
	created.Verified = false

	u.byEmail[email] = created

	return &created, nil
}
