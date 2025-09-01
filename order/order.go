package order

import "food_delivery_service/common"

type Order struct {
    ID           string         `json:"id"`
    RestaurantID string         `json:"restaurant_id"`
    Items        []string       `json:"items"`
    UserLocation common.Location `json:"user_location"`
    DriverID     string         `json:"driver_id"`
    Status       string         `json:"status"`
}

// Init sets up order management
func Init() {}
