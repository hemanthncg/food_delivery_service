package apis

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
)

func TestHealthHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/health", nil)
    w := httptest.NewRecorder()
    HealthHandler(w, req)
    if w.Result().StatusCode != http.StatusOK {
        t.Fatalf("expected 200")
    }
}

// Local integration using miniredis via REDIS_ADDR
func setupFakeRedis(t *testing.T) func() {
    t.Helper()
    addr := os.Getenv("REDIS_ADDR")
    if addr == "" {
        os.Setenv("REDIS_ADDR", "127.0.0.1:6379")
    }
    // Note: we rely on external dev redis or docker
    // for simplicity in this suite.
    return func() {}
}

func TestRestaurantEnrollAndList(t *testing.T) {
    teardown := setupFakeRedis(t)
    defer teardown()
    // enroll
    body := []byte(`{"id":"restT","name":"Test","location":{"lat":1,"lng":2},"menu":["A","B"]}`)
    req := httptest.NewRequest("POST", "/restaurant/enroll", bytes.NewReader(body))
    w := httptest.NewRecorder()
    RestaurantEnrollHandler(w, req)
    if w.Result().StatusCode != http.StatusOK {
        t.Fatalf("enroll failed: %d", w.Result().StatusCode)
    }
    // list
    req2 := httptest.NewRequest("GET", "/restaurants", nil)
    w2 := httptest.NewRecorder()
    RestaurantsHandler(w2, req2)
    if w2.Result().StatusCode != http.StatusOK {
        t.Fatalf("list failed: %d", w2.Result().StatusCode)
    }
}

func TestMenuGet(t *testing.T) {
    teardown := setupFakeRedis(t)
    defer teardown()
    // need a restaurant with menu
    enroll := []byte(`{"id":"restM","name":"Test","location":{"lat":1,"lng":2},"menu":["X"]}`)
    req := httptest.NewRequest("POST", "/restaurant/enroll", bytes.NewReader(enroll))
    w := httptest.NewRecorder()
    RestaurantEnrollHandler(w, req)
    if w.Result().StatusCode != http.StatusOK {
        t.Fatalf("enroll failed: %d", w.Result().StatusCode)
    }
    // fetch menu
    req2 := httptest.NewRequest("GET", "/menu?restaurant_id=restM", nil)
    w2 := httptest.NewRecorder()
    MenuHandler(w2, req2)
    if w2.Result().StatusCode != http.StatusOK {
        t.Fatalf("menu failed: %d", w2.Result().StatusCode)
    }
}

func TestDriverEnrollListAndLocation(t *testing.T) {
    teardown := setupFakeRedis(t)
    defer teardown()
    // enroll driver
    body := []byte(`{"id":"driverT","name":"Bob"}`)
    req := httptest.NewRequest("POST", "/driver/enroll", bytes.NewReader(body))
    w := httptest.NewRecorder()
    DriverEnrollHandler(w, req)
    if w.Result().StatusCode != http.StatusOK {
        t.Fatalf("enroll driver failed")
    }
    // list drivers
    req2 := httptest.NewRequest("GET", "/drivers", nil)
    w2 := httptest.NewRecorder()
    DriversHandler(w2, req2)
    if w2.Result().StatusCode != http.StatusOK {
        t.Fatalf("list drivers failed")
    }
    // update location
    loc := map[string]interface{}{
        "id": "driverT",
        "location": map[string]float64{"lat": 1, "lng": 2},
    }
    b, _ := json.Marshal(loc)
    req3 := httptest.NewRequest("POST", "/driver/location", bytes.NewReader(b))
    w3 := httptest.NewRecorder()
    DriverLocationHandler(w3, req3)
    if w3.Result().StatusCode != http.StatusOK {
        t.Fatalf("update location failed")
    }
    // get location
    req4 := httptest.NewRequest("GET", "/driver/location?id=driverT", nil)
    w4 := httptest.NewRecorder()
    DriverLocationHandler(w4, req4)
    if w4.Result().StatusCode != http.StatusOK {
        t.Fatalf("get location failed")
    }
}
