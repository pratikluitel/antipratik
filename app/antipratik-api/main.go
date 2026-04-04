package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/pratikluitel/antipratik/api"
	"github.com/pratikluitel/antipratik/config"
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

	db, err := store.Open(cfg.DB.Path)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := store.RunMigrations(db); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	if err := store.SeedIfEmpty(db); err != nil {
		log.Fatalf("seed: %v", err)
	}

	postStore := store.NewPostStore(db)
	linkStore := store.NewLinkStore(db)
	postLogic := logic.NewPostService(postStore)
	linkLogic := logic.NewLinkService(linkStore)
	postH := api.NewPostHandler(postLogic)
	linkH := api.NewLinkHandler(linkLogic)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/posts/{slug}", postH.GetPost)
	mux.HandleFunc("GET /api/posts", postH.GetPosts)
	mux.HandleFunc("GET /api/links/featured", linkH.GetFeaturedLinks)
	mux.HandleFunc("GET /api/links", linkH.GetLinks)

	handler := api.CORSMiddleware(mux)

	addr := cfg.Addr()
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
