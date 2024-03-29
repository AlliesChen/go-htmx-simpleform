package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Film struct {
	Title    string
	Director string
}

var (
	//go:embed templates
	templates embed.FS
	//go:embed static
	static embed.FS
	pages  = map[string]string{
		"/":          "templates/index.html",
		"/add-film/": "templates/index.html",
	}
)

func main() {
	fmt.Println("start serving...")

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	fileServer := http.FileServer(http.FS(static))
	router.Handle("/static/*", fileServer)

	// GET /
	getPage := func(w http.ResponseWriter, r *http.Request) {
		page, ok := pages[r.RequestURI]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		tmpl, err := template.ParseFS(templates, page)
		if err != nil {
			log.Printf("page %s not found in pages caches...", r.RequestURI)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		films := map[string][]Film{
			"Films": {
				{Title: "The Godfather", Director: "Francis Ford Coppola"},
				{Title: "Blade Runner", Director: "Ridley Scott"},
				{Title: "The Thing", Director: "John Carpenter"},
			},
		}
		if err := tmpl.Execute(w, films); err != nil {
			log.Printf("template %s execute error: %s", r.RequestURI, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	// POST /add-film
	addFilm := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		page, ok := pages[r.RequestURI]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		tmpl, err := template.ParseFS(templates, page)
		if err != nil {
			log.Printf("page %s not found in pages caches...", r.RequestURI)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		title := r.PostFormValue("title")
		director := r.PostFormValue("director")
		if err := tmpl.ExecuteTemplate(w, "film-list-element", Film{
			Title:    title,
			Director: director,
		}); err != nil {
			log.Printf("template %s execute error: %s", r.RequestURI, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	// Router and handlers
	router.Get("/", getPage)
	router.Post("/add-film/", addFilm)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	defer server.Shutdown(nil)
}
