package usecase

import "github.com/yertaypert/go-assignment7/internal/entity"

type (
	UserRepository interface {
		RegisterUser(user *entity.User) (*entity.User, error)
		LoginUser(user *entity.LoginUserDTO) (*entity.User, error)
	}

	UserInterface interface {
		LoginUser(user *entity.LoginUserDTO) (string, error)
		RegisterUser(user *entity.User) (*entity.User, string, error)
	}
)
