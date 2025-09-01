package menu

import (
    "encoding/json"
    "food_delivery_service/common"
)

const (
    restaurantKey = "restaurants"
    menuPrefix    = "menu:"
)

func AddRestaurantToStore(r Restaurant) error {
    client := common.GetRedisClient()
    data, _ := json.Marshal(r)
    return client.HSet(common.Ctx, restaurantKey, r.ID, data).Err()
}

func GetAllRestaurantsFromStore() ([]Restaurant, error) {
    client := common.GetRedisClient()
    vals, err := client.HVals(common.Ctx, restaurantKey).Result()
    if err != nil {
        return nil, err
    }
    var res []Restaurant
    for _, v := range vals {
        var r Restaurant
        if err := json.Unmarshal([]byte(v), &r); err == nil {
            res = append(res, r)
        }
    }
    return res, nil
}

func AddMenuToStore(restaurantID string, items []string) error {
    client := common.GetRedisClient()
    data, _ := json.Marshal(items)
    return client.Set(common.Ctx, menuPrefix+restaurantID, data, 0).Err()
}

func GetMenuFromStore(restaurantID string) ([]string, error) {
    client := common.GetRedisClient()
    val, err := client.Get(common.Ctx, menuPrefix+restaurantID).Result()
    if err != nil {
        return nil, err
    }
    var items []string
    if err := json.Unmarshal([]byte(val), &items); err != nil {
        return nil, err
    }
    return items, nil
}
