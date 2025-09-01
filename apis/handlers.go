package apis

import (
    "encoding/json"
    "net/http"
    "food_delivery_service/menu"
    "food_delivery_service/driver"
    "food_delivery_service/common"
    "food_delivery_service/order"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func RestaurantEnrollHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", 405)
        return
    }
    var req menu.EnrollRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil || req.ID == "" || req.Name == "" || len(req.Menu) == 0 {
        http.Error(w, "invalid request", 400)
        return
    }
    if err := menu.EnrollRestaurant(req); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func RestaurantsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "method not allowed", 405)
        return
    }
    rest, err := menu.ListRestaurants()
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(rest)
}

func MenuHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("restaurant_id")
    if id == "" {
        http.Error(w, "missing id", 400)
        return
    }
    items, err := menu.GetMenuService(id)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(items)
}

func DriverEnrollHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", 405)
        return
    }
    var req driver.Driver
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil || req.ID == "" || req.Name == "" {
        http.Error(w, "invalid request", 400)
        return
    }
    if err := driver.EnrollDriver(req); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func DriversHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "method not allowed", 405)
        return
    }
    drivers, err := driver.ListDrivers()
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    type out struct {
        ID        string          `json:"id"`
        Name      string          `json:"name"`
        Location  common.Location `json:"location"`
        Available bool            `json:"available"`
    }
    resp := make([]out, 0, len(drivers))
    for _, d := range drivers {
        avail, _ := driver.IsDriverAvailable(d.ID)
        resp = append(resp, out{
            ID: d.ID,
            Name: d.Name,
            Location: d.Location,
            Available: avail,
        })
    }
    json.NewEncoder(w).Encode(resp)
}

func DriverLocationHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        var req struct {
            ID       string         `json:"id"`
            Location common.Location `json:"location"`
        }
        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil || req.ID == "" {
            http.Error(w, "invalid request", 400)
            return
        }
        if err := driver.UpdateLocation(req.ID, req.Location); err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
        return
    }
    if r.Method == http.MethodGet {
        id := r.URL.Query().Get("id")
        if id == "" {
            http.Error(w, "missing id", 400)
            return
        }
        loc, err := driver.GetLocation(id)
        if err != nil {
            http.Error(w, err.Error(), 404)
            return
        }
        avail, _ := driver.IsDriverAvailable(id)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "id":   id,
            "location": loc,
            "available": avail,
        })
        return
    }
    http.Error(w, "method not allowed", 405)
}

func OrderHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        var req order.Order
        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil || req.ID == "" || req.RestaurantID == "" || len(req.Items) == 0 {
            http.Error(w, "invalid request", 400)
            return
        }
        req.Status = "created"
        if err := order.PlaceOrder(&req); err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        json.NewEncoder(w).Encode(req)
        return
    }
    http.Error(w, "method not allowed", 405)
}

func DeliveredHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        var req struct {
            OrderID string `json:"order_id"`
        }
        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil || req.OrderID == "" {
            http.Error(w, "invalid request", 400)
            return
        }
        statusMsg, _ := json.Marshal(map[string]string{
            "order_id": req.OrderID,
            "status":   "Delivered",
        })
        _ = common.PublishKafka("order_status_update", []byte(req.OrderID), statusMsg)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
        return
    }
    http.Error(w, "method not allowed", 405)
}

func GetOrderHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "missing id", 400)
        return
    }
    o, err := order.GetOrder(id)
    if err != nil {
        http.Error(w, err.Error(), 404)
        return
    }
    json.NewEncoder(w).Encode(o)
}

func ListOrdersHandler(w http.ResponseWriter, r *http.Request) {
    orders, err := order.ListOrders()
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(orders)
}
