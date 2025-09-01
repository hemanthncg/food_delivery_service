package order

import (
    "encoding/json"
    "food_delivery_service/common"
    "food_delivery_service/driver"
)

func PlaceOrder(o *Order) error {
    if err := AddOrderToStore(o); err != nil {
        return err
    }
    // Try immediate reservation to avoid races
    // or missing consumer.
    driverID, err := driver.FindAndReserveNearestDriver(
        o.UserLocation, o.ID,
    )
    if err == nil && driverID != "" {
        o.DriverID = driverID
        o.Status = "assigned"
        _ = UpdateOrderInStore(o)
        statusMsg, _ := json.Marshal(map[string]string{
            "order_id": o.ID,
            "status":   "Order assigned to delivery",
        })
        _ = common.PublishKafka(
            "order_status_update",
            []byte(o.ID), statusMsg,
        )
        return nil
    }
    // Fallback to async assignment
    data, _ := json.Marshal(o)
    _ = common.PublishKafka(
        "delivery_assignment",
        []byte(o.ID), data,
    )
    statusMsg, _ := json.Marshal(map[string]string{
        "order_id": o.ID,
        "status":   "Order placed",
    })
    _ = common.PublishKafka(
        "order_status_update",
        []byte(o.ID), statusMsg,
    )
    return nil
}

func AssignDriver(orderID, driverID string) error {
    o, err := GetOrderFromStore(orderID)
    if err != nil {
        return err
    }
    o.DriverID = driverID
    o.Status = "assigned"
    return UpdateOrderInStore(o)
}

func MarkDelivered(orderID string) error {
    o, err := GetOrderFromStore(orderID)
    if err != nil {
        return err
    }
    o.Status = "delivered"
    if o.DriverID != "" {
        _ = driver.ReleaseDriver(o.DriverID)
    }
    return UpdateOrderInStore(o)
}

func GetOrder(orderID string) (*Order, error) {
    return GetOrderFromStore(orderID)
}

func ListOrders() ([]*Order, error) {
    return ListOrdersFromStore()
}

func AssignNearestDriver(orderID string, userLoc common.Location) error {
    driverID, err := driver.FindAndReserveNearestDriver(
        userLoc,
        orderID,
    )
    if err != nil {
        return err
    }
    return AssignDriver(orderID, driverID)
}
