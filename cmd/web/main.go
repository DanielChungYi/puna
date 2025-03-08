package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DanielChungYi/puna/internal/config"
	"github.com/DanielChungYi/puna/internal/driver"
	"github.com/DanielChungYi/puna/internal/handlers"
	"github.com/DanielChungYi/puna/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main function
func main() {
	// change this to true when in production
	app.InProduction = false

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database....")
	dsn := "host=localhost user=postgres password=qwer1234 dbname=testcy port=5432 sslmode=disable"
	_, err := driver.ConnectSQL(dsn)
	if err != nil {
		log.Println("Fail to connect")
	} else {
		log.Println("Success to connect")
	}
	//db.GORM.Close()

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
