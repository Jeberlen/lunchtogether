package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Jeberlen/lunchtogether/graph/model"
	"github.com/Jeberlen/lunchtogether/menu_items"
	"github.com/Jeberlen/lunchtogether/restaurant_menues"
	restaurant "github.com/Jeberlen/lunchtogether/restaurants"
)

// CreateRestaurant is the resolver for the createRestaurant field.
func (r *mutationResolver) CreateRestaurant(ctx context.Context, input model.NewRestaurant) (*model.Restaurant, error) {
	var restaurant restaurant.Restaurant

	restaurant.Name = input.Name
	restaurant.Date = input.Date
	restaurantID := restaurant.Save()

	for _, item := range input.Menu {
		var menu_item menu_items.MenuItem

		menu_item.Type = item.Type
		menu_item.Name = item.Name
		menu_item.Description = item.Description
		menu_item.URL = item.URL
		menu_item.DayOfWeek = item.DayOfWeek

		menu_itemID := menu_item.Save()

		var resturant_menues restaurant_menues.RestaurantMenu
		resturant_menues.RestaurantID = restaurantID
		resturant_menues.MenuItemID = menu_itemID
		resturant_menues.Save()
	}

	return &model.Restaurant{
		ID:   strconv.FormatInt(restaurantID, 10),
		Name: restaurant.Name,
		Date: restaurant.Date,
	}, nil
}

// CreateMenuItem is the resolver for the createMenuItem field.
func (r *mutationResolver) CreateMenuItem(ctx context.Context, input model.NewMenuItem) (*model.MenuItem, error) {
	var menu_item menu_items.MenuItem
	menu_item.Type = input.Type
	menu_item.Name = input.Name
	menu_item.Description = input.Description
	menu_item.URL = input.URL
	menu_item.DayOfWeek = input.DayOfWeek
	menu_itemID := menu_item.Save()

	return &model.MenuItem{
		ID:          strconv.FormatInt(menu_itemID, 10),
		Type:        menu_item.Type,
		Name:        menu_item.Name,
		Description: menu_item.Description,
		URL:         menu_item.URL,
		DayOfWeek:   menu_item.DayOfWeek,
	}, nil
}

// Restaurants is the resolver for the restaurants field.
func (r *queryResolver) Restaurants(ctx context.Context) ([]*model.Restaurant, error) {
	var results []*model.Restaurant
	var dbRestaurant = restaurant.GetAll()

	for _, restaurant := range dbRestaurant {

		var menu_map []*model.MenuItem
		for _, menu := range restaurant.Menu {
			item := &model.MenuItem{
				ID:          menu.ID,
				Type:        menu.Type,
				Name:        menu.Name,
				Description: menu.Description,
				URL:         menu.URL,
				DayOfWeek:   menu.DayOfWeek,
			}
			menu_map = append(menu_map, item)
		}

		results = append(results, &model.Restaurant{ID: restaurant.ID, Name: restaurant.Name, Date: restaurant.Date, Menu: menu_map})
	}
	return results, nil
}

// RestaurantsByDate is the resolver for the restaurantsByDate field.
func (r *queryResolver) RestaurantsByDate(ctx context.Context, date string) ([]*model.Restaurant, error) {
	const shortForm = "2006-Jan-02"
	var dateObject, _ = time.Parse(shortForm, date)
	_, dateAsWeek := dateObject.ISOWeek()

	var results []*model.Restaurant
	dateString := strconv.Itoa(dateAsWeek)
	var dbRestaurant = restaurant.GetResturantByDate(dateString)

	for _, restaurant := range dbRestaurant {
		var menu_map []*model.MenuItem
		for _, menu := range restaurant.Menu {
			item := &model.MenuItem{
				ID:          menu.ID,
				Type:        menu.Type,
				Name:        menu.Name,
				Description: menu.Description,
				URL:         menu.URL,
				DayOfWeek:   menu.DayOfWeek}
			menu_map = append(menu_map, item)
		}

		results = append(results, &model.Restaurant{ID: restaurant.ID, Name: restaurant.Name, Date: restaurant.Date, Menu: menu_map})
	}
	return results, nil
}

// RestaurantsByDateAndDayOfWeek is the resolver for the restaurantsByDateAndDayOfWeek field.
func (r *queryResolver) RestaurantsByDateAndDayOfWeek(ctx context.Context, date string) ([]*model.Restaurant, error) {
	panic(fmt.Errorf("not implemented: RestaurantsByDateAndDayOfWeek - restaurantsByDateAndDayOfWeek"))
}

// TypedMenuByDate is the resolver for the typedMenuByDate field.
func (r *queryResolver) TypedMenuByDate(ctx context.Context, date string) ([]*model.TypedMenuItem, error) {
	var results []*model.TypedMenuItem
	listOfTypes := []string{"meat", "fish", "vegetarian", "other", "salad"}

	for _, typ := range listOfTypes {
		dbRestaurant := menu_items.GetMenuByDateAndType(date, typ)
		var menuItems []*model.MenuItem
		for _, menu := range dbRestaurant {
			restaurantByMenuId := restaurant.GetRestaurantByMenuId(menu.ID)
			restaurant := &model.Restaurant{
				ID:   restaurantByMenuId.ID,
				Name: restaurantByMenuId.Name,
				Date: restaurantByMenuId.Date,
			}

			menuItem := &model.MenuItem{
				ID:          menu.ID,
				Name:        menu.Name,
				Description: menu.Description,
				Type:        menu.Type,
				URL:         menu.URL,
				DayOfWeek:   menu.DayOfWeek,
				Restaurant:  restaurant,
			}
			menuItems = append(menuItems, menuItem)
		}

		var typedMenuItem model.TypedMenuItem
		typedMenuItem.Type = typ
		typedMenuItem.Menu = menuItems

		results = append(results, &typedMenuItem)
	}

	return results, nil
}

// TypedMenuByDateAndType is the resolver for the typedMenuByDateAndType field.
func (r *queryResolver) TypedMenuByDateAndType(ctx context.Context, date string, typeArg string) ([]*model.TypedMenuItem, error) {
	var results []*model.TypedMenuItem

	dbRestaurant := menu_items.GetMenuByDateAndType(date, typeArg)
	var menuItems []*model.MenuItem
	for _, menu := range dbRestaurant {
		restaurantByMenuId := restaurant.GetRestaurantByMenuId(menu.ID)
		restaurant := &model.Restaurant{
			ID:   restaurantByMenuId.ID,
			Name: restaurantByMenuId.Name,
			Date: restaurantByMenuId.Date,
		}

		menuItem := &model.MenuItem{
			ID:          menu.ID,
			Name:        menu.Name,
			Description: menu.Description,
			Type:        menu.Type,
			URL:         menu.URL,
			DayOfWeek:   menu.DayOfWeek,
			Restaurant:  restaurant,
		}
		menuItems = append(menuItems, menuItem)
	}

	var typedMenuItem model.TypedMenuItem
	typedMenuItem.Type = typeArg
	typedMenuItem.Menu = menuItems

	results = append(results, &typedMenuItem)

	return results, nil
}

// MenuItemsByType is the resolver for the menuItemsByType field.
func (r *queryResolver) MenuItemsByType(ctx context.Context, typeArg string) ([]*model.MenuItem, error) {
	panic(fmt.Errorf("not implemented: MenuItemsByType - menuItemsByType"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func filter(objects []*model.MenuItem, condition func(*model.MenuItem) bool) []*model.MenuItem {
	var filtered []*model.MenuItem
	for _, obj := range objects {
		if condition(obj) {
			filtered = append(filtered, obj)
		}
	}
	return filtered
}
func getUnique[T any, Key comparable](objects []*T, getKey func(*T) Key) []*T {
	// Create a map to store unique objects based on the key attribute
	uniqueObjectsMap := make(map[Key]*T)

	// Iterate over the objects, adding unique objects to the map
	for _, obj := range objects {
		key := getKey(obj)
		if _, ok := uniqueObjectsMap[key]; !ok {
			uniqueObjectsMap[key] = obj
		}
	}

	// Convert the unique objects map back to a slice
	var uniqueObjects []*T
	for _, obj := range uniqueObjectsMap {
		uniqueObjects = append(uniqueObjects, obj)
	}

	return uniqueObjects
}
