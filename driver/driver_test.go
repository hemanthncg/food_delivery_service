package driver

import (
    "food_delivery_service/common"
    "testing"
)

func ensureRedis(t *testing.T) bool {
    c := common.GetRedisClient()
    if err := c.Ping(common.Ctx).Err(); err != nil {
        t.Log("redis not available:", err)
        return false
    }
    return true
}

func TestUpdateLocation(t *testing.T) {
    if !ensureRedis(t) {
        t.Skip("redis not running")
    }
}

func TestFindAndReserveNearestDriver(t *testing.T) {
    if !ensureRedis(t) {
        t.Skip("redis not running")
    }
    c := common.GetRedisClient()
    _ = c.Del(common.Ctx, "drivers").Err()
    _ = c.Del(common.Ctx, "drivers_geo").Err()

    d1 := Driver{ID: "d1", Name: "A"}
    d2 := Driver{ID: "d2", Name: "B"}
    d3 := Driver{ID: "d3", Name: "C"}

    if err := AddDriverToStore(d1); err != nil {
        t.Fatal(err)
    }
    if err := AddDriverToStore(d2); err != nil {
        t.Fatal(err)
    }
    if err := AddDriverToStore(d3); err != nil {
        t.Fatal(err)
    }

    // user near (19.0, 72.9)
    u := common.Location{Lat: 19.0, Lng: 72.9}
    // place drivers around
    _ = UpdateDriverLocationInStore(
        "d1", common.Location{Lat: 19.01, Lng: 72.90},
    )
    _ = UpdateDriverLocationInStore(
        "d2", common.Location{Lat: 19.10, Lng: 72.90},
    )
    _ = UpdateDriverLocationInStore(
        "d3", common.Location{Lat: 19.02, Lng: 72.91},
    )

    // first pick should be d1
    id, err := FindAndReserveNearestDriver(u, "o1")
    if err != nil {
        t.Fatal(err)
    }
    if id != "d1" {
        t.Fatalf("want d1 got %s", id)
    }
    // second pick should skip d1
    id2, err := FindAndReserveNearestDriver(u, "o2")
    if err != nil {
        t.Fatal(err)
    }
    if id2 != "d3" {
        t.Fatalf("want d3 got %s", id2)
    }
    // release d1 and pick again
    if err := ReleaseDriver("d1"); err != nil {
        t.Fatal(err)
    }
    id3, err := FindAndReserveNearestDriver(u, "o3")
    if err != nil {
        t.Fatal(err)
    }
    if id3 != "d1" {
        t.Fatalf("want d1 got %s", id3)
    }
}
