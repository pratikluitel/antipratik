package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/pratikluitel/antipratik/api"
	"github.com/pratikluitel/antipratik/common/logging"
	"github.com/pratikluitel/antipratik/config"
	"github.com/pratikluitel/antipratik/db"
	"github.com/pratikluitel/antipratik/logic"
	"github.com/pratikluitel/antipratik/store"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger := logging.New(cfg.Logging.Level)

	sqlDB, err := db.Open(cfg.DB.Path)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer func() {
		if err = sqlDB.Close(); err != nil {
			log.Printf("error closing db: %v", err)
		}
	}()

	if err = db.RunMigrations(sqlDB); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	fileStore, err := store.NewFileStore(cfg.Storage)
	if err != nil {
		log.Fatalf("init file store: %v", err)
	}

	postStore := store.NewPostStore(sqlDB)
	linkStore := store.NewLinkStore(sqlDB)
	userStore := store.NewUserStore(sqlDB)
	newsletterStore := store.NewNewsletterStore(sqlDB)
	settingsStore := store.NewSettingsStore(sqlDB)

	setupSvc := logic.NewSetupService(userStore, settingsStore)

	ctx := context.Background()

	if cfg.AdminPassword != "" {
		if err = setupSvc.UpsertAdminUser(ctx, cfg.AdminPassword); err != nil {
			log.Fatalf("upsert admin user: %v", err)
		}
	}

	jwtSecret, err := setupSvc.GetOrCreateJWTSecret(ctx)
	if err != nil {
		log.Fatalf("jwt secret: %v", err)
	}

	uploadSvc := logic.NewUploadService(fileStore)
	postLogic := logic.NewPostService(postStore, fileStore, logger)
	linkLogic := logic.NewLinkService(linkStore)
	authService := logic.NewAuthService(userStore, jwtSecret)
	newsletterLogic := logic.NewNewsletterService(newsletterStore)

	postH := api.NewPostHandler(postLogic, uploadSvc, logger)
	linkH := api.NewLinkHandler(linkLogic, logger)
	authH := api.NewAuthHandler(authService, logger)
	fileH := api.NewFileServingHandler(fileStore, logger)
	newsletterH := api.NewNewsletterHandler(newsletterLogic, logger)

	mux := http.NewServeMux()
	RegisterRoutes(mux, postH, linkH, authH, authService, fileH, newsletterH, "api/openapi.yaml", "api/swagger.html")

	handler := api.CORSMiddleware(mux)

	addr := cfg.Addr()
	logger.Info("listening on", "addr", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
