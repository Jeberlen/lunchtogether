package hivecrawler

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Jeberlen/lunchtogether/menu_items"
	"github.com/Jeberlen/lunchtogether/restaurants"
	"github.com/gocolly/colly"
)

var collector *colly.Collector
var meatTypes = []string{
	"ham",
	"entrecôte", "turkey", "lamb",
	"veal", "rib", "fläsk", "pork",
	"beef", "chicken", "pannbiff",
	"bacon", "sirloin", "steak", "burger",
	"brisket", "kebab"}
var fishTypes = []string{
	"fish", "shrimp", "prawn", "clam",
	"clams", "crab", "salmon", "tuna",
	"tilapia", "cod", "snapper", "sardine",
	"sardines", "herring", "haddock",
	"flounder", "trout", "pollock", "bass",
	"halibut", "pike", "mackerel"}
var vegTypes = []string{"vegetarian", "vegan", "veggie", "halloumi", "pannoumi"}

func RemoveStrings(inputStrings []string, toRemove ...string) []string {
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

func Contains(slice []string, str string) bool {

	for _, item := range slice {
		if item == str {
			return true
		}
	}

	return false
}

func GetDayOfWeek(day string) (string, string) {
	switch day {
	case "MONDAY":
		return "1", "https://thehivefoodmarket.se/#monday"
	case "TUESDAY":
		return "2", "https://thehivefoodmarket.se/#tuesday"
	case "WEDNESDAY":
		return "3", "https://thehivefoodmarket.se/#wednesday"
	case "THURSDAY":
		return "4", "https://thehivefoodmarket.se/#thursday"
	case "FRIDAY":
		return "5", "https://thehivefoodmarket.se/#friday"
	default:
		return "1", "https://thehivefoodmarket.se/"
	}
}

func GetFoodTypeFromName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	slicesOfName := strings.Split(name, " ")

	for _, slice := range slicesOfName {

		wordInName := strings.TrimSpace(slice)
		wordInName = strings.ToLower(wordInName)

		if Contains(fishTypes, wordInName) {
			return "fish"
		} else if Contains(meatTypes, wordInName) {
			return "meat"
		} else if Contains(vegTypes, wordInName) {
			return "vegetarian"
		}
	}
	return "other"
}

func GetFoodTypeFromDescription(description string) string {
	description = strings.ToLower(description)
	description = strings.ReplaceAll(description, ",", "")
	sliceOfDesc := strings.Split(description, " ")

	for _, slice := range sliceOfDesc {
		desc := strings.TrimSpace(slice)
		desc = strings.ToLower(desc)

		if Contains(meatTypes, desc) {
			return "meat"
		} else if Contains(fishTypes, desc) {
			return "fish"
		} else if Contains(vegTypes, desc) {
			return "vegetarian"
		}
	}
	return "other"
}

func InitSpider() {
	collector = colly.NewCollector()
}

func StartCrawl(waitGroup *sync.WaitGroup) {

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

		var salads []menu_items.MenuItem
		h.ForEach("#av_section_1", func(i int, h *colly.HTMLElement) {
			h.ForEach(".OYPEnA", func(i int, h *colly.HTMLElement) {
				var menuItem menu_items.MenuItem
				salladSlice := strings.Split(h.Text, "\n")
				name := salladSlice[0]
				desc := salladSlice[1]

				menuItem.Name = name
				menuItem.Description = desc
				menuItem.Type = "salad"
				menuItem.URL = "https://thehivefoodmarket.se/"

				for i := 1; i < 6; i++ {
					menuItem.DayOfWeek = strconv.Itoa(i)
					salads = append(salads, menuItem)
				}
			})
		})

		menuItems = append(menuItems, salads...)
		h.ForEach(".avia-section", func(i int, h *colly.HTMLElement) {
			day := h.ChildText(".av-special-heading")
			h.ForEachWithBreak(".av_textblock_section", func(i int, h *colly.HTMLElement) bool {
				dailySlice := strings.Split(h.Text, "\n")
				stringsToRemove := []string{"", "MONDO", "SPICE CLUB", "HUSMANSKOST", "WEST COAST", "PIZZA", "TRUE FOOD", "Salads Of The Week"}
				filteredStrings := RemoveStrings(dailySlice, stringsToRemove...)
				if len(filteredStrings) == 0 {
					return false
				}
				for i, food := range filteredStrings {
					var hiveMenuItem menu_items.MenuItem
					if i%2 == 0 {

						day, url := GetDayOfWeek(day)
						hiveMenuItem.DayOfWeek = day
						hiveMenuItem.URL = url
						hiveMenuItem.Name = food
						hiveMenuItem.Description = filteredStrings[i+1]

						foodType := GetFoodTypeFromName(hiveMenuItem.Name)
						if foodType == "other" {
							foodType = GetFoodTypeFromDescription(hiveMenuItem.Description)
						}
						hiveMenuItem.Type = foodType

						menuItems = append(menuItems, hiveMenuItem)
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
	waitGroup.Done()

}
