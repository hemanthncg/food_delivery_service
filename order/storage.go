package order

import (
    "encoding/json"
    "food_delivery_service/common"
)

const (
    orderKey = "orders"
)

func AddOrderToStore(o *Order) error {
    client := common.GetRedisClient()
    data, _ := json.Marshal(o)
    return client.HSet(common.Ctx, orderKey, o.ID, data).Err()
}

func GetOrderFromStore(id string) (*Order, error) {
    client := common.GetRedisClient()
    val, err := client.HGet(common.Ctx, orderKey, id).Result()
    if err != nil {
        return nil, err
    }
    var o Order
    if err := json.Unmarshal([]byte(val), &o); err != nil {
        return nil, err
    }
    return &o, nil
}

func UpdateOrderInStore(o *Order) error {
    return AddOrderToStore(o)
}

func ListOrdersFromStore() ([]*Order, error) {
    client := common.GetRedisClient()
    vals, err := client.HVals(common.Ctx, orderKey).Result()
    if err != nil {
        return nil, err
    }
    var res []*Order
    for _, v := range vals {
        var o Order
        if err := json.Unmarshal([]byte(v), &o); err == nil {
            res = append(res, &o)
        }
    }
    return res, nil
}
