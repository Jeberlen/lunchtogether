package harvestcrawler

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

	url := "https://harvestrestaurant.se/#lunch"

	InitSpider()

	collector.OnRequest(func(r *colly.Request) {
		log.Print("Visiting: ", r.URL)
	})

	collector.OnHTML("#main", func(h *colly.HTMLElement) {
		var restaurant restaurants.Restaurant
		restaurant.Name = "Harvest by Mannerström"
		_, currentWeek := time.Now().ISOWeek()
		restaurant.Date = strconv.Itoa(currentWeek)
		var menuItems []menu_items.MenuItem

		h.ForEach(".content", func(i int, h *colly.HTMLElement) {
			log.Print(h.Text)
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
