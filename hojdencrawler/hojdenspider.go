package hojdencrawler

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/Jeberlen/lunchtogether/menu_items"
	"github.com/Jeberlen/lunchtogether/restaurants"
	"github.com/gocolly/colly"
)

var collector *colly.Collector

func KeepEveryThird(list []menu_items.MenuItem) []menu_items.MenuItem {
	var result []menu_items.MenuItem
	for i := 0; i < len(list); i += 3 {
		result = append(result, list[i])
	}
	return result
}

func GetDayOfWeek(day string) string {
	switch day {
	case "Måndag":
		return "1"
	case "Tisdag":
		return "2"
	case "Onsdag":
		return "3"
	case "Torsdag":
		return "4"
	case "Fredag":
		return "5"
	}
	return "1"
}

func HasTypePrefix(text string) bool {
	if strings.HasPrefix(text, "GOOD VEGETARIAN") ||
		strings.HasPrefix(text, "MEAT") ||
		strings.HasPrefix(text, "FISH") ||
		strings.HasPrefix(text, "STREET") ||
		strings.HasPrefix(text, "FUSION") {
		return true
	}

	return false
}

func InitSpider() {
	collector = colly.NewCollector()
}

func StartCrawl(waitGroup *sync.WaitGroup) {

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

				menuItem.DayOfWeek = GetDayOfWeek(strings.TrimSpace(splitForDay[1]))

				h.ForEach(".eng-meny", func(i int, h *colly.HTMLElement) {
					if HasTypePrefix(h.Text) {
						slices := strings.Split(h.Text, ": ")
						for _, slice := range slices {
							if strings.HasPrefix(slice, "GOOD VEGETARIAN") {
								menuItem.Type = "vegetarian"
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

		menuItems = KeepEveryThird(menuItems)
		var ptrMenuItems []*menu_items.MenuItem
		for i := range menuItems {
			menuItemPointer := &menuItems[i]
			ptrMenuItems = append(ptrMenuItems, menuItemPointer)
		}

		restaurant.Menu = ptrMenuItems

		log.Print("Starting to save complete restaurant")

		restaurant.SaveCompleteRestaurant()
		waitGroup.Done()

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

	waitGroup.Done()

}
