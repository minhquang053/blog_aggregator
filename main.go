package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/minhquang053/blog_aggregator/internal/database"
)

func main() {
	godotenv.Load()
	db, err := sql.Open("postgres", os.Getenv("DBURL"))
	if err != nil {
		log.Println(err)
		return
	}

	apiCfg := apiConfig{
		DB: database.New(db),
	}

	r := chi.NewRouter()

	// Define middlewares
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Register endpoints with v1Router
	v1Router := chi.NewRouter()
	v1Router.Get("/readiness", handlerReadiness)
	v1Router.Get("/err", handlerError)
	v1Router.Post("/users", apiCfg.handlerUsersCreate)

	// Mount subrouters to main router
	r.Mount("/v1", v1Router)

	server := http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}

	log.Fatal(server.ListenAndServe())
}
