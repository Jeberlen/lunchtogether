package restaurants

import (
	"fmt"
	"log"

	database "github.com/Jeberlen/lunchtogether/db"
	menu_items "github.com/Jeberlen/lunchtogether/menu_items"
)

type Restaurant struct {
	ID		string
	Name    string
  	Date    string
  	Menu	[]*menu_items.MenuItem
}

func (restaurant Restaurant) Save() int64 {
	var id int64
	query := fmt.Sprintf("INSERT INTO restaurant(name,date) VALUES('%s','%s') RETURNING id", 
		restaurant.Name, 
		restaurant.Date)

	err := database.Db.QueryRow(query).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Row inserted.")

	return id
}

func GetAll() []Restaurant {
	stmt, err := database.Db.Prepare("select id, name, date from restaurant")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var restaurants []Restaurant
	for rows.Next() {
		var restaurant Restaurant
		err := rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Date)
		if err != nil{
			log.Fatal(err)
		}

		query := fmt.Sprintf("select id, type, name, description, url from menu_item where id=%s", restaurant.ID)
		stmt, err := database.Db.Prepare(query)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		rows, err := stmt.Query()
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var menu []*menu_items.MenuItem
		for rows.Next() {
			var menu_item menu_items.MenuItem
			err := rows.Scan(
				&menu_item.ID, 
				&menu_item.Type, 
				&menu_item.Name, 
				&menu_item.Description, 
				&menu_item.URL)
			if err != nil{
				log.Fatal(err)
			}
			menu = append(menu, &menu_item)
		}

		restaurant.Menu = menu

		restaurants = append(restaurants, restaurant)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return restaurants
}

func GetResturantByDate(date string) []Restaurant {
	stmt, err := database.Db.Prepare("select id, name, date from restaurant")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var restaurants []Restaurant
	for rows.Next() {
		var restaurant Restaurant
		err := rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Date)
		if err != nil{
			log.Fatal(err)
		}

		query := fmt.Sprintf("select id, type, name, description, url from menu_item where id=%s", restaurant.ID)
		stmt, err := database.Db.Prepare(query)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		rows, err := stmt.Query()
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var menu []*menu_items.MenuItem
		for rows.Next() {
			var menu_item menu_items.MenuItem
			err := rows.Scan(
				&menu_item.ID, 
				&menu_item.Type, 
				&menu_item.Name, 
				&menu_item.Description, 
				&menu_item.URL)
			if err != nil{
				log.Fatal(err)
			}
			menu = append(menu, &menu_item)
		}

		restaurant.Menu = menu

		restaurants = append(restaurants, restaurant)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return restaurants
}