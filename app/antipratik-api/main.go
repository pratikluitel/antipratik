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
	"github.com/pratikluitel/antipratik/components/broadcaster/lib/resend"
	broadcasterlogic "github.com/pratikluitel/antipratik/components/broadcaster/logic"
	broadcastersvc "github.com/pratikluitel/antipratik/components/broadcaster/services"
	broadcasterstore "github.com/pratikluitel/antipratik/components/broadcaster/store"
	filesapi "github.com/pratikluitel/antipratik/components/files/api"
	fileslogic "github.com/pratikluitel/antipratik/components/files/logic"
	filesservices "github.com/pratikluitel/antipratik/components/files/services"
	filesstore "github.com/pratikluitel/antipratik/components/files/store"
	postsapi "github.com/pratikluitel/antipratik/components/posts/api"
	postslogic "github.com/pratikluitel/antipratik/components/posts/logic"
	postsservices "github.com/pratikluitel/antipratik/components/posts/services"
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
	uploadLogic := fileslogic.NewUploadLogic(fileStore)
	storageSvc := filesservices.NewStorageService(fileStore)
	uploaderSvc := filesservices.NewUploaderService(uploadLogic)
	fileH := filesapi.NewFileServingHandler(fileStore, logger)

	// Auth component
	userStore := authstore.NewUserStore(sqlDB)
	settingsStore := authstore.NewSettingsStore(sqlDB)
	setupLogic := authlogic.NewSetupLogic(userStore, settingsStore)

	ctx := context.Background()

	if cfg.AdminPassword != "" {
		if err = setupLogic.UpsertAdminUser(ctx, cfg.AdminPassword); err != nil {
			log.Fatalf("upsert admin user: %v", err)
		}
	}

	jwtSecret, err := setupLogic.GetOrCreateJWTSecret(ctx)
	if err != nil {
		log.Fatalf("jwt secret: %v", err)
	}

	authLogic := authlogic.NewAuthLogic(userStore, jwtSecret)
	authH := authapi.NewAuthHandler(authLogic, logger)

	// Posts component
	postStore := postsstore.NewPostStore(sqlDB)
	linkStore := postsstore.NewLinkStore(sqlDB)

	postLogic := postslogic.NewPostLogic(postStore, storageSvc, logger)
	linkLogic := postslogic.NewLinkLogic(linkStore)

	postsSvc := postsservices.NewPostsService(postLogic)

	postH := postsapi.NewPostHandler(postLogic, uploaderSvc, logger)
	linkH := postsapi.NewLinkHandler(linkLogic, logger)

	// Broadcaster component
	broadcasterStore := broadcasterstore.NewBroadcasterStore(sqlDB)

	resendClient := resend.NewClient(resend.Config{
		APIKey:   cfg.Broadcaster.Resend.APIKey,
		Host:     cfg.Broadcaster.Resend.Host,
		Port:     cfg.Broadcaster.Resend.Port,
		From:     cfg.Broadcaster.Resend.FromEmail,
		FromName: cfg.Broadcaster.Resend.FromName,
	}, logger)

	broadcasterLogic, err := broadcasterlogic.NewBroadcasterLogic(
		broadcasterStore,
		resendClient,
		postsSvc,
		cfg.AdminEmail,
		cfg.SiteDomain,
		"antipratik",
		cfg.Broadcaster.Resend.FromName,
		logger,
	)
	if err != nil {
		log.Fatalf("init broadcaster: %v", err)
	}

	_ = broadcastersvc.NewSubscriberService(broadcasterLogic)
	broadcasterH := broadcasterapi.NewBroadcasterHandler(broadcasterLogic, logger)

	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, postH, linkH, authH, authLogic, fileH, broadcasterH, "swagger/openapi.yaml", "swagger/swagger.html")

	handler := handlers.CORSMiddleware(mux)

	addr := cfg.Addr()
	logger.Info("listening on", "addr", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
