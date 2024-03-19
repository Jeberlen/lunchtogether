package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
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
				DayOfWeek:   menu.DayOfWeek}
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
	dateString := strconv.Itoa(dateAsWeek + 1)
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

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }