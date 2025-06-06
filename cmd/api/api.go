package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sikozonpc/social/internal/store"
	//httpSwagger "github.com/swaggo/http-swagger/v2"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr   string
	db     dbConfig
	env    string
	apiURL string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	fmt.Println("***A")
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
	fmt.Println("***B")
	r.Route("/v1", func(r chi.Router) {
		fmt.Println("***C")
		r.Get("/health", app.healthCheckHandler)

		//docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		//r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))
		fmt.Println("***D")
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			fmt.Println("***E")
			r.Route("/{postID}", func(r chi.Router) {
				fmt.Println("***F")
				//r.Use(app.postsContextMiddleware)

				r.Get("/", app.getPostHandler)
				fmt.Println("***G")
				r.Delete("/", app.deletePostHandler)
				fmt.Println("***H")
				r.Patch("/", app.updatePostHandler)
				fmt.Println("***I")
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)

				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

	})
	fmt.Println("ASD")
	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	//docs.SwaggerInfo.Version = version
	//docs.SwaggerInfo.Host = app.config.apiURL
	//docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}
