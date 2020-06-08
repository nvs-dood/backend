package main

//go run github.com/99designs/gqlgen generate

import (
	"log"
	"net/http"
	"os"

	"github.com/EnglederLucas/nvs-dood/auth"
	"github.com/EnglederLucas/nvs-dood/database"
	"github.com/EnglederLucas/nvs-dood/graph"
	"github.com/EnglederLucas/nvs-dood/graph/generated"
	"github.com/go-chi/chi"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "3000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(auth.Middleware())

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		DB: database.GetDB(),
	}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
