package main

import (
	"context"
	"embed"
	"flag"
	"log"
	"net/http"

	authapi "github.com/pratikluitel/antipratik/components/auth/api"
	authlogic "github.com/pratikluitel/antipratik/components/auth/logic"
	authstore "github.com/pratikluitel/antipratik/components/auth/store"
	broadcasterapi "github.com/pratikluitel/antipratik/components/broadcaster/api"
	broadcasterlogic "github.com/pratikluitel/antipratik/components/broadcaster/logic"
	broadcasterstore "github.com/pratikluitel/antipratik/components/broadcaster/store"
	broadcastersvc "github.com/pratikluitel/antipratik/components/broadcaster/services"
	filesapi "github.com/pratikluitel/antipratik/components/files/api"
	fileslogic "github.com/pratikluitel/antipratik/components/files/logic"
	filesservices "github.com/pratikluitel/antipratik/components/files/services"
	filesstore "github.com/pratikluitel/antipratik/components/files/store"
	postsapi "github.com/pratikluitel/antipratik/components/posts/api"
	postslogic "github.com/pratikluitel/antipratik/components/posts/logic"
	postsstore "github.com/pratikluitel/antipratik/components/posts/store"
	"github.com/pratikluitel/antipratik/common/db"
	"github.com/pratikluitel/antipratik/common/logging"
	"github.com/pratikluitel/antipratik/config"
	"github.com/pratikluitel/antipratik/handlers"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

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

	if err = db.RunMigrations(sqlDB, migrationsFS); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	// Files component
	fileStore, err := filesstore.NewFileStore(cfg.Storage)
	if err != nil {
		log.Fatalf("init file store: %v", err)
	}
	uploadSvc := fileslogic.NewUploadService(fileStore)
	storageSvc := filesservices.NewStorageService(fileStore)
	uploaderSvc := filesservices.NewUploaderService(uploadSvc)
	fileH := filesapi.NewFileServingHandler(fileStore, logger)

	// Auth component
	userStore := authstore.NewUserStore(sqlDB)
	settingsStore := authstore.NewSettingsStore(sqlDB)
	setupSvc := authlogic.NewSetupService(userStore, settingsStore)

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

	authService := authlogic.NewAuthService(userStore, jwtSecret)
	authH := authapi.NewAuthHandler(authService, logger)

	// Posts component
	postStore := postsstore.NewPostStore(sqlDB)
	linkStore := postsstore.NewLinkStore(sqlDB)

	postLogic := postslogic.NewPostService(postStore, storageSvc, logger)
	linkLogic := postslogic.NewLinkService(linkStore)

	postH := postsapi.NewPostHandler(postLogic, uploaderSvc, logger)
	linkH := postsapi.NewLinkHandler(linkLogic, logger)

	// Broadcaster component
	nlStore := broadcasterstore.NewNewsletterStore(sqlDB)
	nlLogic := broadcasterlogic.NewNewsletterService(nlStore)
	_ = broadcastersvc.NewSubscriberService(nlLogic)
	newsletterH := broadcasterapi.NewNewsletterHandler(nlLogic, logger)

	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, postH, linkH, authH, authService, fileH, newsletterH, "swagger/openapi.yaml", "swagger/swagger.html")

	handler := handlers.CORSMiddleware(mux)

	addr := cfg.Addr()
	logger.Info("listening on", "addr", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
