package order

import (
    "encoding/json"
    "food_delivery_service/common"
    "food_delivery_service/driver"
    "log"
)

func StartConsumers() {
    go deliveryAssignmentConsumer()
    go orderStatusUpdateConsumer()
}

func deliveryAssignmentConsumer() {
    common.ConsumeKafka("delivery_assignment", "delivery-group", func(msg []byte) {
        var o Order
        if err := json.Unmarshal(msg, &o); err != nil {
            log.Println("delivery assign decode error:", err)
            return
        }
        driverID, err := driver.FindAndReserveNearestDriver(
            o.UserLocation,
            o.ID,
        )
        if err == nil && driverID != "" {
            o.DriverID = driverID
            o.Status = "assigned"
            _ = UpdateOrderInStore(&o)
            statusMsg, _ := json.Marshal(map[string]string{
                "order_id": o.ID,
                "status":   "Order assigned to delivery",
            })
            _ = common.PublishKafka("order_status_update", []byte(o.ID), statusMsg)
        }
    })
}

func orderStatusUpdateConsumer() {
    common.ConsumeKafka("order_status_update", "status-group", func(msg []byte) {
        var s struct {
            OrderID string `json:"order_id"`
            Status  string `json:"status"`
        }
        if err := json.Unmarshal(msg, &s); err != nil {
            log.Println("status update decode error:", err)
            return
        }
        o, err := GetOrderFromStore(s.OrderID)
        if err == nil && o != nil {
            o.Status = s.Status
            _ = UpdateOrderInStore(o)
            if s.Status == "Delivered" && o.DriverID != "" {
                _ = driver.ReleaseDriver(o.DriverID)
            }
        }
    })
}
