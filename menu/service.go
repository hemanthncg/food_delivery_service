package menu

import "food_delivery_service/common"

type EnrollRequest struct {
    ID       string         `json:"id"`
    Name     string         `json:"name"`
    Location common.Location `json:"location"`
    Menu     []string       `json:"menu"`
}

func EnrollRestaurant(req EnrollRequest) error {
    r := Restaurant{ID: req.ID, Name: req.Name, Location: req.Location}
    if err := AddRestaurantToStore(r); err != nil {
        return err
    }
    return AddMenuToStore(req.ID, req.Menu)
}

func ListRestaurants() ([]Restaurant, error) {
    return GetAllRestaurantsFromStore()
}

func GetMenuService(restaurantID string) ([]string, error) {
    return GetMenuFromStore(restaurantID)
}
