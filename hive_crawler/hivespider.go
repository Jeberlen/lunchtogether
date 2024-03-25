package crawler

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Jeberlen/lunchtogether/menu_items"
	"github.com/Jeberlen/lunchtogether/restaurants"
	"github.com/gocolly/colly"
)

var collector *colly.Collector

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func allFieldsNotEmpty(s interface{}) bool {
	v := reflect.ValueOf(s)

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)

		if reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(fieldValue.Type()).Interface()) {
			return false
		}
	}

	return true
}

func removeStrings(inputStrings []string, toRemove ...string) []string {
	result := make([]string, 0)
	removeSet := make(map[string]bool)

	// Create a set for strings to remove
	for _, s := range toRemove {
		removeSet[s] = true
	}

	// Iterate through the input strings, appending non-removed strings to the result
	for _, s := range inputStrings {
		s := strings.TrimSpace(s)
		if !removeSet[s] {
			result = append(result, s)
		}
	}

	return result
}

func contains(slice []string, str string) bool {

	for _, item := range slice {
		if item == str {
			return true
		}
	}

	return false
}

func InitSpider() {
	collector = colly.NewCollector()
}

func StartCrawl() {

	url := "https://thehivefoodmarket.se/"

	InitSpider()

	collector.OnRequest(func(r *colly.Request) {
		log.Print("Visiting: ", r.URL)
	})

	collector.OnHTML("#main", func(h *colly.HTMLElement) {
		var restaurant restaurants.Restaurant

		restaurant.Name = "The Hive"
		_, currentWeek := time.Now().ISOWeek()
		restaurant.Date = strconv.Itoa(currentWeek)

		var menuItems []menu_items.MenuItem
		h.ForEach(".avia-section", func(i int, h *colly.HTMLElement) {
			day := h.ChildText(".av-special-heading")
			h.ForEachWithBreak(".av_textblock_section", func(i int, h *colly.HTMLElement) bool {
				dailySlice := strings.Split(h.Text, "\n")
				type HiveMenuItem struct {
					Name        string
					Description string
					Type        string
					URL         string
					DayOfWeek   string
				}

				stringsToRemove := []string{"", "MONDO", "SPICE CLUB", "HUSMANSKOST", "WEST COAST", "PIZZA", "TRUE FOOD", "Salads Of The Week"}
				filteredStrings := removeStrings(dailySlice, stringsToRemove...)
				if len(filteredStrings) == 0 {
					return false
				}
				for i, food := range filteredStrings {
					var hiveMenuItem HiveMenuItem
					if i%2 == 0 {
						switch day {
						case "MONDAY":
							hiveMenuItem.DayOfWeek = "1"
							hiveMenuItem.URL = "https://thehivefoodmarket.se/#monday"
						case "TUESDAY":
							hiveMenuItem.DayOfWeek = "2"
							hiveMenuItem.URL = "https://thehivefoodmarket.se/#tuesday"
						case "WEDNESDAY":
							hiveMenuItem.DayOfWeek = "3"
							hiveMenuItem.URL = "https://thehivefoodmarket.se/#wednesday"
						case "THURSDAY":
							hiveMenuItem.DayOfWeek = "4"
							hiveMenuItem.URL = "https://thehivefoodmarket.se/#thursday"
						case "FRIDAY":
							hiveMenuItem.DayOfWeek = "5"
							hiveMenuItem.URL = "https://thehivefoodmarket.se/#friday"
						default:
							hiveMenuItem.DayOfWeek = "1"
							hiveMenuItem.URL = "https://thehivefoodmarket.se/"
						}

						hiveMenuItem.Name = food
						hiveMenuItem.Description = filteredStrings[i+1]

						lowerCaseForComp := strings.ToLower(hiveMenuItem.Description)
						noUnusedChars := strings.ReplaceAll(lowerCaseForComp, ",", "")
						sliceOfDesc := strings.Split(noUnusedChars, " ")

						meatTypes := []string{"fläsk", "pork", "beef", "chicken", "ham", "pannbiff", "bacon", "sirloin", "steak", "burger", "brisket"}
						fishTypes := []string{"fish", "shrimp"}
						saladTypes := []string{"salad"}
						pizzaTypes := []string{"pizza"}
						vegTypes := []string{"vegetarian", "vegan"}

						for _, slice := range sliceOfDesc {
							name := strings.TrimSpace(hiveMenuItem.Name)
							name = strings.ToLower(name)
							desc := strings.TrimSpace(slice)
							desc = strings.ToLower(desc)

							if contains(meatTypes, name) || contains(meatTypes, desc) {
								hiveMenuItem.Type = "meat"
								break // Stop after first match
							} else if contains(fishTypes, name) || contains(fishTypes, desc) {
								hiveMenuItem.Type = "fish"
								break // Stop after first match
							} else if contains(saladTypes, name) || contains(saladTypes, desc) {
								hiveMenuItem.Type = "salad"
								break // Stop after first match
							} else if contains(pizzaTypes, name) || contains(pizzaTypes, desc) {
								hiveMenuItem.Type = "pizza"
								break // Stop after first match
							} else if contains(vegTypes, name) || contains(vegTypes, desc) {
								hiveMenuItem.Type = "vegetarian"
								break // Stop after first match
							} else {
								hiveMenuItem.Type = "other"
							}
						}

						actualMenuItem := menu_items.MenuItem{
							Name:        hiveMenuItem.Name,
							Description: hiveMenuItem.Description,
							Type:        hiveMenuItem.Type,
							URL:         hiveMenuItem.URL,
							DayOfWeek:   hiveMenuItem.DayOfWeek,
						}

						menuItems = append(menuItems, actualMenuItem)

					}
				}
				return false
			})
		})

		var ptrMenuItems []*menu_items.MenuItem
		for i := range menuItems {
			menuItemPointer := &menuItems[i]
			ptrMenuItems = append(ptrMenuItems, menuItemPointer)
		}

		restaurant.Menu = ptrMenuItems

		log.Print(restaurant.Name)
		log.Print(restaurant.Date)
		for _, e := range restaurant.Menu {
			log.Print(e.Name)
			log.Print(e.Description)
			log.Print(e.Type)
			log.Print(e.DayOfWeek)
			log.Print("-------------------------")

		}

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
