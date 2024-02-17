package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	//movies
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission("movies:read", app.listMoviesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission("movies:write", app.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission("movies:read", app.showMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission("movies:write", app.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovieHandler))
	//users
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	//users auth
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/jwtlogin", app.loginWithJwtTokenHandler)
	//metrics
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	//
	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}

// package main

// import (
// 	"net/http"

// 	"github.com/go-chi/chi"
// 	"github.com/go-chi/cors"
// )

// func (app *application) routes() http.Handler {
// 	router := chi.NewRouter()

// 	// Basic CORS
// 	router.Use(cors.Handler(cors.Options{
// 		AllowedOrigins:   []string{"https://*", "http://*"},
// 		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
// 		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
// 		ExposedHeaders:   []string{"Link"},
// 		AllowCredentials: false,
// 		MaxAge:           300,
// 	}))
// 	// router := httprouter.New()
// 	v1Router := chi.NewRouter()

// 	// router.NotFound = http.HandlerFunc(app.notFoundResponse)
// 	// router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
// 	router.MethodNotAllowed(http.HandlerFunc(app.methodNotAllowedResponse))
// 	router.NotFound(http.HandlerFunc(app.notFoundResponse))

// 	// router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
// 	v1Router.Get("/healthcheck", app.healthcheckHandler)
// 	//movies
// 	// router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission("movies:read", app.listMoviesHandler))
// 	v1Router.Get("/movies", app.requirePermission("movies:read", app.listMoviesHandler))
// 	// router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission("movies:write", app.createMovieHandler))
// 	v1Router.Post("/movies", app.requirePermission("movies:write", app.createMovieHandler))
// 	// router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission("movies:read", app.showMovieHandler))
// 	v1Router.Get("/movies/{id}", app.requirePermission("movies:read", app.showMovieHandler))
// 	// router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission("movies:write", app.updateMovieHandler))
// 	v1Router.Patch("/movies/{id}", app.requirePermission("movies:write", app.updateMovieHandler))
// 	// router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovieHandler))
// 	v1Router.Delete("/movies/{id}", app.requirePermission("movies:write", app.deleteMovieHandler))
// 	//users
// 	// router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
// 	v1Router.Post("/users", app.registerUserHandler)
// 	//users auth
// 	// router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
// 	v1Router.Put("/users/activated", app.activateUserHandler)
// 	// router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
// 	v1Router.Post("/tokens/authentication", app.createAuthenticationTokenHandler)
// 	//metrics
// 	// router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())
// 	// v1Router.Get("/debug/vars", expvar.Handler())
// 	//
// 	router.Mount("/v1", v1Router)
// 	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
// }
