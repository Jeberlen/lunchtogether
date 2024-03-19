package menu_items

import (
	"fmt"
	"log"

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
