package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Jeberlen/lunchtogether/crawler"
	database "github.com/Jeberlen/lunchtogether/db"
	"github.com/Jeberlen/lunchtogether/graph"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func crawlerHandler(w http.ResponseWriter, r *http.Request) {
	crawler.StartCrawl()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	database.InitDB()
	defer database.CloseDB()

	// Create a GraphQL server
	gqlHandler := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{},
	}))

	// Define CORS options to allow all origins
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := cors.Handler(gqlHandler)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", handler)
	http.HandleFunc("/crawl", crawlerHandler)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
