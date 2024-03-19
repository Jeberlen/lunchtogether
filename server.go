package main

import (
	"log"
	"net/http"
	"os"

	crawler "github.com/Jeberlen/lunchtogether/hojdencrawler"
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

	log.Print("starting to crawl")
	crawler.StartCrawl()
	log.Print("ending crawl")
	//database.InitDB()
	//defer database.CloseDB()

	//srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//http.Handle("/query", srv)
	//http.HandleFunc("/crawl", crawlerHandler)

	//log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	//log.Fatal(http.ListenAndServe(":"+port, nil))
}
