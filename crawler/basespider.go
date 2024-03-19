package crawler

import (
	"fmt"
	"log"

	"github.com/Jeberlen/lunchtogether/restaurants"
	"github.com/gocolly/colly"
)

func StartCrawl() {

	collector := colly.NewCollector()
	collector.Visit("https://volvo-cars.nordrest.se/hojden/")

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	collector.OnHTML(".menu-col", func(e *colly.HTMLElement) {
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

}
