package app

import (
	"fmt"
	"os"

	httpcontroller "github.com/yertaypert/go-assignment7/internal/controller/http"
	"github.com/yertaypert/go-assignment7/internal/usecase"
	"github.com/yertaypert/go-assignment7/internal/usecase/repo"
	"github.com/yertaypert/go-assignment7/pkg"
)

func Run() error {
	pg, err := pkg.NewPostgres()
	if err != nil {
		return fmt.Errorf("init postgres: %w", err)
	}
	defer pg.Close()

	userRepo := repo.NewUserRepo(pg)
	userUC := usecase.NewUserUseCase(userRepo)
	router := httpcontroller.NewRouter(userUC)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return router.Run(fmt.Sprintf(":%s", port))
}
