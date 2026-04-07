package app

import (
	"fmt"
	"os"

	httpcontroller "github.com/yertaypert/go-assignment7/internal/controller/http"
	"github.com/yertaypert/go-assignment7/internal/usecase"
	"github.com/yertaypert/go-assignment7/internal/usecase/repo"
)

func Run() error {
	userRepo := repo.NewUserRepo()
	userUC := usecase.NewUserUseCase(userRepo)
	router := httpcontroller.NewRouter(userUC)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return router.Run(fmt.Sprintf(":%s", port))
}
