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

	mux.Get("/search-availability", handlers.Repo.Availability)     // 顯示畫面
	mux.Get("/check-availability", handlers.Repo.CheckAvailability) // 查可用數量
	mux.Post("/make-reservation", handlers.Repo.MakeReservation)    // 確定預約

	mux.Get("/court-info", handlers.Repo.CourtInfo)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/login", handlers.Repo.Login)
	mux.Post("/login", handlers.Repo.PostLogin)
	mux.Get("/logout", handlers.Repo.Logout)
	mux.Get("/register", handlers.Repo.Register)
	mux.Post("/register", handlers.Repo.PostRegister)
	const (
		assetsDir     = "./static"
		assetsUrlPath = "/assets"
	)
	fileServer := http.FileServer(http.Dir(assetsDir))
	mux.Handle(assetsUrlPath+"/*", http.StripPrefix(assetsUrlPath, fileServer))

	return mux
}
