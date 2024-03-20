package crawler

import (
	"fmt"
	"log"
	"reflect"
	"slices"
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

// getUnique returns unique objects based on a specified key attribute
func getUnique[T any, Key comparable](objects []T, getKey func(T) Key) []T {
	// Create a map to store unique objects based on the key attribute
	uniqueObjectsMap := make(map[Key]T)

	// Iterate over the objects, adding unique objects to the map
	for _, obj := range objects {
		key := getKey(obj)
		if _, ok := uniqueObjectsMap[key]; !ok {
			uniqueObjectsMap[key] = obj
		}
	}

	// Convert the unique objects map back to a slice
	var uniqueObjects []T
	for _, obj := range uniqueObjectsMap {
		uniqueObjects = append(uniqueObjects, obj)
	}

	return uniqueObjects
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
		var menuItems []menu_items.MenuItem

		restaurant.Name = "The Hive"
		_, currentWeek := time.Now().ISOWeek()
		restaurant.Date = strconv.Itoa(currentWeek)

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

		uniqueSalads := getUnique(salads, func(obj menu_items.MenuItem) string {
			return obj.Name
		})

		menuItems = append(menuItems, uniqueSalads...)
		var dailyFood []menu_items.MenuItem
		meatTypes := []string{"pork", "beef", "chicken", "ham", "pannbiff", "bacon"}
		fishTypes := []string{"fish", "shrimp"}

		h.ForEach(".avia-section", func(i int, h *colly.HTMLElement) {
			day := h.ChildText(".av-special-heading")
			h.ForEach(".av_textblock_section", func(i int, h *colly.HTMLElement) {
				dailySlice := strings.Split(h.Text, "\n")
				type HiveMenuItem struct {
					Name        string
					Description string
					Type        string
					URL         string
					DayOfWeek   string
				}
				var hiveMenuItem HiveMenuItem
				for i, food := range dailySlice {

					if i < len(dailySlice)-1 {
						if len(food) == 0 {
							remove(dailySlice, i+1)
							continue
						}

						switch food {
						case "MONDO":
							remove(dailySlice, i+1)
							continue
						case "SPICE CLUB":
							remove(dailySlice, i+1)
							continue
						case "HUSMANSKOST":
							remove(dailySlice, i+1)
							continue
						case "W EST COAST":
							remove(dailySlice, i+1)
							continue
						case "PIZZA":
							remove(dailySlice, i+1)
							continue
						case "TRUE FOOD":
							remove(dailySlice, i+1)
							continue
						}
					}

				}

				for _, food := range dailySlice[1:] {
					switch day {
					case "MONDAY":
						hiveMenuItem.DayOfWeek = "1"
					case "TUESDAY":
						hiveMenuItem.DayOfWeek = "2"
					case "WEDNESDAY":
						hiveMenuItem.DayOfWeek = "3"
					case "THURSDAY":
						hiveMenuItem.DayOfWeek = "4"
					case "FRIDAY":
						hiveMenuItem.DayOfWeek = "5"
					}

					hiveMenuItem.URL = "https://thehivefoodmarket.se/#monday"

					if len(food) < 45 && len(food) != 0 {
						hiveMenuItem.Name = food
					} else if len(food) >= 45 {

						hiveMenuItem.Description = food
						lowerCaseForComp := strings.ToLower(hiveMenuItem.Description)
						noUnusedChars := strings.ReplaceAll(lowerCaseForComp, ",", "")
						sliceOfDesc := strings.Split(noUnusedChars, " ")

						for _, slice := range sliceOfDesc {
							if slices.Contains(meatTypes, slice) || slices.Contains(meatTypes, strings.ToLower(hiveMenuItem.Name)) {
								hiveMenuItem.Type = "meat"
								break
							} else if slices.Contains(fishTypes, slice) {
								hiveMenuItem.Type = "fish"
								break
							} else {
								hiveMenuItem.Type = "vegitarian"
							}
						}
					}

					if allFieldsNotEmpty(hiveMenuItem) {
						actualMenuItem := menu_items.MenuItem{
							Name:        hiveMenuItem.Name,
							Description: hiveMenuItem.Description,
							Type:        hiveMenuItem.Type,
							URL:         hiveMenuItem.URL,
							DayOfWeek:   hiveMenuItem.DayOfWeek,
						}
						dailyFood = append(dailyFood, actualMenuItem)
					}
				}
			})
		})

		uniqueMenuItems := getUnique(dailyFood, func(obj menu_items.MenuItem) string {
			return obj.Name
		})

		menuItems = append(menuItems, uniqueMenuItems...)
		var ptrMenuItems []*menu_items.MenuItem
		for i := range menuItems {
			menuItemPointer := &menuItems[i]
			ptrMenuItems = append(ptrMenuItems, menuItemPointer)
		}

		restaurant.Menu = ptrMenuItems
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
