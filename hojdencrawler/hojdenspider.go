package crawler

import (
	"fmt"
	"log"

	"github.com/Jeberlen/lunchtogether/restaurants"
	"github.com/gocolly/colly"
)

var collector *colly.Collector

func InitSpider() {
	collector = colly.NewCollector()
}

func StartCrawl() {

	InitSpider()

	collector.OnRequest(func(r *colly.Request) {
		log.Print("Visiting: ", r.URL)
	})

	collector.OnHTML(".menu-col", func(e *colly.HTMLElement) {
		log.Print(e)
		var restaurantList []restaurants.Restaurant
		restaurant := restaurants.Restaurant{}

		// scraping the data of interest
		restaurant.Name = e.ChildAttr("a", "href")
		restaurant.Date = e.ChildAttr("img", "src")
		//restaurant.Menu = e.ChildText("h2") TODO: Get menuitems

		// adding the product instance with scraped data to the list of restaurants
		restaurantList = append(restaurantList, restaurant)

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

	log.Print("Sending spider to https://volvo-cars.nordrest.se/hojden/")
	err := collector.Visit("https://volvo-cars.nordrest.se/hojden/")
	if err != nil {
		log.Fatal(err)
	}

}
