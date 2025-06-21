package main

import (
	"net/http"

	"github.com/DanielChungYi/puna/internal/config"
	"github.com/DanielChungYi/puna/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Get("/make-reservation", handlers.Repo.Reservation)

	mux.Get("/court-info", handlers.Repo.CourtInfo)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/login", handlers.Repo.ShowLogin)
	mux.Post("/login", handlers.Repo.PostShowLogin)
	mux.Get("/register", handlers.Repo.ShowRegister)
	mux.Post("/register", handlers.Repo.PostShowRegister)
	const (
		assetsDir     = "./static"
		assetsUrlPath = "/assets"
	)
	fileServer := http.FileServer(http.Dir(assetsDir))
	mux.Handle(assetsUrlPath+"/*", http.StripPrefix(assetsUrlPath, fileServer))

	return mux
}
