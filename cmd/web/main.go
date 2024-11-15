package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/karahmon/hotel-bookings/pkg/config"
	"github.com/karahmon/hotel-bookings/pkg/handlers"
	"github.com/karahmon/hotel-bookings/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":3000"

var app config.Appconfig
var session *scs.SessionManager

func main() {

	// change this to true in production
	app.InProduction = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = app.InProduction
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		fmt.Println("Cannot Create Template Cache")
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println("Starting Application on port", portNumber)
	// _ = http.ListenAndServe(portNumber, nil)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
