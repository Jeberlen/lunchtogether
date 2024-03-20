package menu_items

import (
	"fmt"
	"log"
	"time"

	database "github.com/Jeberlen/lunchtogether/db"
)

type MenuItem struct {
	ID          string `json:"id"`
	DayOfWeek   string `json:"dayOfWeek"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func (menu_item MenuItem) Save() int64 {
	var id int64
	query := fmt.Sprintf("INSERT INTO menu_item(type,name,description,url,dayOfWeek) VALUES('%s','%s','%s','%s','%s') RETURNING id",
		menu_item.Type,
		menu_item.Name,
		menu_item.Description,
		menu_item.URL,
		menu_item.DayOfWeek)

	err := database.Db.QueryRow(query).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Row inserted.")

	return id
}

func GetMenuByDate(date string) []MenuItem {

	dateAsObject, err := time.Parse("2006-Jan-02", date)
	if err != nil {
		log.Fatal(err)
	}
	dayOfWeek := int(dateAsObject.Weekday())

	stmt, err := database.Db.Prepare(
		fmt.Sprintf("select id, name, description, type, dayOfWeek, url from menu_item where dayOfWeek = '%d'", dayOfWeek),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	menuItemRows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer menuItemRows.Close()

	var menuItems []MenuItem
	for menuItemRows.Next() {
		var menuItem MenuItem
		err := menuItemRows.Scan(
			&menuItem.ID,
			&menuItem.Name,
			&menuItem.Description,
			&menuItem.Type,
			&menuItem.DayOfWeek,
			&menuItem.URL,
		)
		if err != nil {
			log.Fatal(err)
		}

		menuItems = append(menuItems, menuItem)
	}

	return menuItems
}
