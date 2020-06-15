package main

//go run github.com/99designs/gqlgen generate

import (
	"fmt"
	"github.com/EnglederLucas/nvs-dood/auth"
	"github.com/EnglederLucas/nvs-dood/graph/generated"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/EnglederLucas/nvs-dood/database"
	"github.com/EnglederLucas/nvs-dood/graph"
	"github.com/go-chi/chi"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
)

const defaultPort = "3000"

func main() {

	ch := time.After(10 * time.Second)
	defer (func() { fmt.Println("waiting"); <-ch; fmt.Println("waited") })()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8100"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)
	router.Use(auth.Middleware())

	db, err := database.GetDB()
	if err != nil {
		return
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		DB: db,
	}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
