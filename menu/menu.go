package menu

import "food_delivery_service/common"

type Restaurant struct {
    ID       string         `json:"id"`
    Name     string         `json:"name"`
    Location common.Location `json:"location"`
}

type Menu struct {
    RestaurantID string   `json:"restaurant_id"`
    Items        []string `json:"items"`
}

var restaurants = []Restaurant{
    {ID: "rest1", Name: "Pizza Place", Location: common.Location{Lat: 19.1, Lng: 72.9}},
    {ID: "rest2", Name: "Burger Joint", Location: common.Location{Lat: 19.2, Lng: 72.8}},
}

var menus = map[string][]string{
    "rest1": {"Margherita", "Pepperoni"},
    "rest2": {"Cheeseburger", "Veggie Burger"},
}

// Init sets up menu management
func Init() {}

// GetAllRestaurants returns all restaurants
func GetAllRestaurants() []Restaurant {
    return restaurants
}

// GetMenu returns a menu for a restaurant
func GetMenu(restaurantID string) *Menu {
    items := menus[restaurantID]
    return &Menu{RestaurantID: restaurantID, Items: items}
}

// AddRestaurant enrolls a new restaurant
func AddRestaurant(id, name string, loc common.Location, items []string) bool {
    for _, r := range restaurants {
        if r.ID == id {
            return false
        }
    }
    restaurants = append(restaurants, Restaurant{ID: id, Name: name, Location: loc})
    menus[id] = items
    return true
}
