package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	database "github.com/Jeberlen/lunchtogether/db"
	"github.com/Jeberlen/lunchtogether/graph"
	harvestCrawler "github.com/Jeberlen/lunchtogether/harvestcrawler"
	hiveCrawler "github.com/Jeberlen/lunchtogether/hivecrawler"
	hojdenCrawler "github.com/Jeberlen/lunchtogether/hojdencrawler"

	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	database.InitDB()
	defer database.CloseDB()

	var waitGroup sync.WaitGroup
	waitGroup.Add(4)

	log.Print("starting to crawl")
	go hojdenCrawler.StartCrawl(&waitGroup)
	go hiveCrawler.StartCrawl(&waitGroup)
	go harvestCrawler.StartCrawl(&waitGroup)
	waitGroup.Wait()

	log.Print("ending crawl")

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

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
