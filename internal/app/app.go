package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	httpcontroller "github.com/yertaypert/go-assignment7/internal/controller/http"
	"github.com/yertaypert/go-assignment7/internal/usecase"
	"github.com/yertaypert/go-assignment7/internal/usecase/repo"
	"github.com/yertaypert/go-assignment7/pkg"
)

func Run() error {
	if err := godotenv.Load(); err != nil {
		// We don't necessarily return err here because production
		// environments might use actual env vars instead of a file
		fmt.Println("No .env file found, reading from system environment")
	}
	pg, err := pkg.NewPostgres()
	if err != nil {
		return fmt.Errorf("init postgres: %w", err)
	}
	defer pg.Close()

	userRepo := repo.NewUserRepo(pg)
	if err := ensureDefaultAdmin(userRepo); err != nil {
		return err
	}
	userUC := usecase.NewUserUseCase(userRepo)
	router := httpcontroller.NewRouter(userUC)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return router.Run(fmt.Sprintf(":%s", port))
}

func ensureDefaultAdmin(userRepo *repo.UserRepo) error {
	adminEmail := strings.TrimSpace(os.Getenv("ADMIN_EMAIL"))
	adminPassword := strings.TrimSpace(os.Getenv("ADMIN_PASSWORD"))
	adminUsername := strings.TrimSpace(os.Getenv("ADMIN_USERNAME"))

	if adminEmail == "" && adminPassword == "" {
		return nil
	}
	if adminEmail == "" || adminPassword == "" {
		return fmt.Errorf("ADMIN_EMAIL and ADMIN_PASSWORD must both be set")
	}

	if _, err := userRepo.EnsureAdminUser(adminUsername, adminEmail, adminPassword); err != nil {
		return fmt.Errorf("ensure default admin: %w", err)
	}

	return nil
}
