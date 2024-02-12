package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Film struct {
	Title    string
	Director string
}

func main() {
	fmt.Println("start serving...")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// GET /
	getPage := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))

		films := map[string][]Film{
			"Films": {
				{Title: "The Godfather", Director: "Francis Ford Coppola"},
				{Title: "Blade Runner", Director: "Ridley Scott"},
				{Title: "The Thing", Director: "John Carpenter"},
			},
		}
		tmpl.Execute(w, films)
	}
	// POST /add-film
	addFilm := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		title := r.PostFormValue("title")
		director := r.PostFormValue("director")
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "film-list-element", Film{
			Title:    title,
			Director: director,
		})
	}
	// Router and handlers
	r.Get("/", getPage)
	r.Post("/add-film/", addFilm)

	log.Fatal(http.ListenAndServe(":8000", r))
}
