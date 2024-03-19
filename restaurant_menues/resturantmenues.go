package restaurant_menues

import (
	"fmt"
	"log"

	database "github.com/Jeberlen/lunchtogether/db"
)


type RestaurantMenu struct {
	RestaurantID		int64
	MenuItemID    		int64
}

func (restaurant_menu RestaurantMenu) Save() {

	query := fmt.Sprintf("INSERT INTO restaurant_menu(restaurant_id,menu_item_id) VALUES(%d,%d)", 
		restaurant_menu.RestaurantID, 
		restaurant_menu.MenuItemID)

	database.Db.QueryRow(query)

	log.Print("Row inserted.")
}