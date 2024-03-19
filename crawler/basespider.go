package crawler

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

type PokemonProduct struct { 
	url, image, name, price string 
}

func StartCrawl() string {

	collector := colly.NewCollector()
	collector.Visit("https://scrapeme.live/shop/")


	collector.OnRequest(func(r *colly.Request) { 
		fmt.Println("Visiting: ", r.URL) 
	}) 

	var ret string

	collector.OnHTML("li.product", func(e *colly.HTMLElement) { 
		var pokemonProducts []PokemonProduct
		pokemonProduct := PokemonProduct{} 
 
		// scraping the data of interest 
		pokemonProduct.url = e.ChildAttr("a", "href") 
		pokemonProduct.image = e.ChildAttr("img", "src") 
		pokemonProduct.name = e.ChildText("h2") 
		pokemonProduct.price = e.ChildText(".price") 
	
		// adding the product instance with scraped data to the list of products 
		pokemonProducts = append(pokemonProducts, pokemonProduct) 

		ret = pokemonProduct.name
	})
	
	collector.OnError(func(_ *colly.Response, err error) { 
		log.Println("Something went wrong: ", err) 
	}) 
	
	collector.OnResponse(func(r *colly.Response) { 
		fmt.Println("Page visited: ", r.Request.URL) 
	}) 
	
	collector.OnScraped(func(r *colly.Response) { 
		fmt.Println(r.Request.URL, " scraped!") 
	})

	return ret
}