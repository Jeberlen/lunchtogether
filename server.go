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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	database.InitDB()
	defer database.CloseDB()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	http.HandleFunc("/crawl", crawlerHandler)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	handler := cors.Default().Handler(srv)

	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatal(err)
	}
}
