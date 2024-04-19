package restaurants

import (
	"fmt"
	"log"

	database "github.com/Jeberlen/lunchtogether/db"
	menu_items "github.com/Jeberlen/lunchtogether/menu_items"
)

type Restaurant struct {
	ID   string
	Name string
	Date string
	Menu []*menu_items.MenuItem
}

func (restaurant Restaurant) SaveCompleteRestaurant() {
	checkIfExistQuery := fmt.Sprintf("SELECT COUNT(*) FROM restaurant WHERE name='%s' AND date='%s'",
		restaurant.Name,
		restaurant.Date)

	rows, err := database.Db.Query(checkIfExistQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var numberOfElements string
	for rows.Next() {
		err := rows.Scan(&numberOfElements)
		if err != nil {
			log.Fatal(err)
		}
	}

	if numberOfElements == "0" {
		query := fmt.Sprintf("INSERT INTO restaurant(name,date) VALUES('%s','%s') RETURNING id",
			restaurant.Name,
			restaurant.Date)

		var restaurantId int64
		err := database.Db.QueryRow(query).Scan(&restaurantId)
		if err != nil {
			log.Fatal(err)
		}

		var menuIds []int64
		for _, item := range restaurant.Menu {
			menuId := item.Save()
			menuIds = append(menuIds, menuId)
		}

		for _, menuId := range menuIds {
			query := fmt.Sprintf("INSERT INTO restaurant_menu(restaurant_id,menu_item_id) VALUES(%d,%d)",
				restaurantId,
				menuId)

			database.Db.QueryRow(query)
		}
		log.Print("Complete restaurant inserted")
	}

	log.Print("Restaurant for date already present")
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

	restaurantRows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer restaurantRows.Close()

	var restaurants []Restaurant
	for restaurantRows.Next() {
		var restaurant Restaurant
		err := restaurantRows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Date)
		if err != nil {
			log.Fatal(err)
		}

		rowsOfMenuItemIds, _ := database.Db.Query(
			fmt.Sprintf("select menu_item_id from restaurant_menu where restaurant_id = %s", restaurant.ID),
		)

		var menu []*menu_items.MenuItem
		for rowsOfMenuItemIds.Next() {
			var menuItemId string
			rowsOfMenuItemIds.Scan(&menuItemId)

			query := fmt.Sprintf("select id, type, name, description, url, dayOfWeek from menu_item where id=%s", menuItemId)
			stmt, err := database.Db.Prepare(query)
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			menuItemRows, err := stmt.Query()
			if err != nil {
				log.Fatal(err)
			}
			defer menuItemRows.Close()

			for menuItemRows.Next() {
				var menu_item menu_items.MenuItem
				err := menuItemRows.Scan(
					&menu_item.ID,
					&menu_item.Type,
					&menu_item.Name,
					&menu_item.Description,
					&menu_item.URL,
					&menu_item.DayOfWeek)
				if err != nil {
					log.Fatal(err)
				}
				menu = append(menu, &menu_item)
			}

			restaurant.Menu = menu

		}
		restaurants = append(restaurants, restaurant)

	}

	if err = restaurantRows.Err(); err != nil {
		log.Fatal(err)
	}

	return restaurants
}

func GetResturantByDate(date string) []Restaurant {
	stmt, err := database.Db.Prepare(
		fmt.Sprintf("select id, name, date from restaurant where date = '%s'", date),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	restaurantRows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer restaurantRows.Close()

	var restaurants []Restaurant
	for restaurantRows.Next() {
		var restaurant Restaurant
		err := restaurantRows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Date)
		if err != nil {
			log.Fatal(err)
		}

		rowsOfMEnuItemIds, _ := database.Db.Query(
			fmt.Sprintf("select menu_item_id from restaurant_menu where restaurant_id = %s", restaurant.ID),
		)

		var menu []*menu_items.MenuItem
		for rowsOfMEnuItemIds.Next() {
			var menuItemId string
			rowsOfMEnuItemIds.Scan(&menuItemId)

			query := fmt.Sprintf("select id, type, name, description, url, dayOfWeek from menu_item where id=%s", menuItemId)
			stmt, err := database.Db.Prepare(query)
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			menuItemRows, err := stmt.Query()
			if err != nil {
				log.Fatal(err)
			}
			defer menuItemRows.Close()

			for menuItemRows.Next() {
				var menu_item menu_items.MenuItem
				err := menuItemRows.Scan(
					&menu_item.ID,
					&menu_item.Type,
					&menu_item.Name,
					&menu_item.Description,
					&menu_item.URL,
					&menu_item.DayOfWeek)
				if err != nil {
					log.Fatal(err)
				}
				menu = append(menu, &menu_item)
			}

			restaurant.Menu = menu

		}
		restaurants = append(restaurants, restaurant)

	}

	if err = restaurantRows.Err(); err != nil {
		log.Fatal(err)
	}

	return restaurants
}

func GetRestaurantByMenuId(id string) *Restaurant {
	stmt, err := database.Db.Prepare(
		fmt.Sprintf("select restaurant_id from restaurant_menu where menu_item_id = '%s'", id),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	restaurantIdRows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}

	var restaurant Restaurant
	for restaurantIdRows.Next() {
		var restaurantId string
		err := restaurantIdRows.Scan(&restaurantId)
		if err != nil {
			log.Fatal(err)
		}

		stmt, err := database.Db.Prepare(
			fmt.Sprintf("select id, name, date from restaurant where id = '%s'", restaurantId),
		)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		restaurants, err := stmt.Query()
		if err != nil {
			log.Fatal(err)
		}

		for restaurants.Next() {
			restaurants.Scan(
				&restaurant.ID,
				&restaurant.Name,
				&restaurant.Date,
			)
		}
	}

	return &restaurant
}

func GetRestaurantByMenuIdAndDate(id string, date string) *Restaurant {
	stmt, err := database.Db.Prepare(
		fmt.Sprintf("select restaurant_id from restaurant_menu where menu_item_id = '%s'", id),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	restaurantIdRows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}

	var restaurant Restaurant
	for restaurantIdRows.Next() {
		var restaurantId string
		err := restaurantIdRows.Scan(&restaurantId)
		if err != nil {
			log.Fatal(err)
		}

		stmt, err := database.Db.Prepare(
			fmt.Sprintf("select id, name, date from restaurant where id = '%s' and date = '%s'", restaurantId, date),
		)
		if err != nil {
			log.Fatal("ERROR IS HERE")

			log.Fatal(err)
		}
		defer stmt.Close()

		restaurants, err := stmt.Query()
		if err != nil {
			log.Fatal(err)
		}

		for restaurants.Next() {
			restaurants.Scan(
				&restaurant.ID,
				&restaurant.Name,
				&restaurant.Date,
			)
		}
	}

	return &restaurant
}
