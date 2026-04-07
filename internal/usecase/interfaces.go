package usecase

import "github.com/yertaypert/go-assignment7/internal/entity"

type (
	UserRepository interface {
		RegisterUser(user *entity.User) (*entity.User, error)
	}

	UserInterface interface {
		RegisterUser(user *entity.User) (*entity.User, string, error)
	}
)
