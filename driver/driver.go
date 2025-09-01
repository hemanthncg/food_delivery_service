package driver

import "food_delivery_service/common"

type Driver struct {
    ID       string         `json:"id"`
    Name     string         `json:"name"`
    Location common.Location `json:"location"`
}

// Init sets up driver management
func Init() {}
