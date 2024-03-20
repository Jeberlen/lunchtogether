package crawler

import (
	"fmt"
	"log"
	"strings"

	"github.com/Jeberlen/lunchtogether/menu_items"
	"github.com/Jeberlen/lunchtogether/restaurants"
	"github.com/gocolly/colly"
)

var collector *colly.Collector

func keepEveryThird(list []menu_items.MenuItem) []menu_items.MenuItem {
	var result []menu_items.MenuItem
	for i := 0; i < len(list); i += 3 {
		result = append(result, list[i])
	}
	return result
}

func InitSpider() {
	collector = colly.NewCollector()
}

func StartCrawl() {

	url := "https://volvo-cars.nordrest.se/hojden/"

	InitSpider()

	collector.OnRequest(func(r *colly.Request) {
		log.Print("Visiting: ", r.URL)
	})

	collector.OnHTML("#current", func(h *colly.HTMLElement) {
		var restaurant restaurants.Restaurant
		var menuItems []menu_items.MenuItem

		h.ForEach(".menu-heading", func(i int, h *colly.HTMLElement) {
			name := h.ChildText("h3")[0:18]
			date := h.ChildText("h2")
			restaurant.Name = name
			restaurant.Date = date[len(date)-2:]
		})

		h.ForEach(".menu-col", func(i int, h *colly.HTMLElement) {

			var menuItem menu_items.MenuItem
			h.ForEach(".menu-item", func(i int, h *colly.HTMLElement) {

				splitForDay := strings.Split(h.Text, "\n")

				switch strings.TrimSpace(splitForDay[1]) {
				case "Måndag":
					menuItem.DayOfWeek = "1"
				case "Tisdag":
					menuItem.DayOfWeek = "2"
				case "Onsdag":
					menuItem.DayOfWeek = "3"
				case "Torsdag":
					menuItem.DayOfWeek = "4"
				case "Fredag":
					menuItem.DayOfWeek = "5"
				}

				h.ForEach(".eng-meny", func(i int, h *colly.HTMLElement) {
					if strings.HasPrefix(h.Text, "GOOD VEGETARIAN") ||
						strings.HasPrefix(h.Text, "MEAT") ||
						strings.HasPrefix(h.Text, "FISH") ||
						strings.HasPrefix(h.Text, "STREET") ||
						strings.HasPrefix(h.Text, "FUSION") {

						slices := strings.Split(h.Text, ": ")
						for _, slice := range slices {
							if strings.HasPrefix(slice, "GOOD VEGETARIAN") {
								menuItem.Type = "vegitarian"
							}
							if strings.HasPrefix(slice, "MEAT") {
								menuItem.Type = "meat"
							}
							if strings.HasPrefix(slice, "FISH") {
								menuItem.Type = "fish"
							}
							if strings.HasPrefix(slice, "STREET") || strings.HasPrefix(slice, "FUSION") {
								menuItem.Type = "other"
							}

							nameAndDesc := strings.Split(slice, " | ")
							for i, part := range nameAndDesc {
								if i == 0 {
									menuItem.Name = part
								}

								nameAndDesc := nameAndDesc[1:]

								desc := strings.Join(nameAndDesc, ", ")
								menuItem.Description = desc
							}
						}
					}

					menuItem.URL = url
					menuItems = append(menuItems, menuItem)

				})

			})
		})

		menuItems = keepEveryThird(menuItems)
		var ptrMenuItems []*menu_items.MenuItem
		for i := range menuItems {
			menuItemPointer := &menuItems[i]
			ptrMenuItems = append(ptrMenuItems, menuItemPointer)
		}

		restaurant.Menu = ptrMenuItems

		log.Print("Starting to save complete restaurant")

		restaurant.SaveCompleteRestaurant()

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

	log.Print("Sending spider to " + url)
	err := collector.Visit(url)
	if err != nil {
		log.Fatal(err)
	}

}
