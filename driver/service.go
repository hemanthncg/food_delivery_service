package driver

import (
    "food_delivery_service/common"
    "github.com/go-redis/redis/v8"
    "time"
)

func EnrollDriver(d Driver) error {
    return AddDriverToStore(d)
}

func ListDrivers() ([]Driver, error) {
    return GetAllDriversFromStore()
}

func UpdateLocation(id string, loc common.Location) error {
    return UpdateDriverLocationInStore(id, loc)
}

func GetLocation(id string) (common.Location, error) {
    d, err := GetDriverFromStore(id)
    if err != nil {
        return common.Location{}, err
    }
    return d.Location, nil
}

func FindNearestDriver(userLoc common.Location) (string, error) {
    // Use Redis GEOSEARCH to find the closest driver.
    client := common.GetRedisClient()
    q := &redis.GeoRadiusQuery{
        Radius: 10000,
        Unit:   "m",
        Sort:   "ASC",
        Count:  1,
        WithDist: true,
    }
    locs, err := client.GeoRadius(
        common.Ctx,
        driversGeoKey,
        userLoc.Lng,
        userLoc.Lat,
        q,
    ).Result()
    if err != nil || len(locs) == 0 {
        return "", err
    }
    return locs[0].Name, nil
}

const (
    driverLockPrefix = "driver_lock:"
)

func ReleaseDriver(driverID string) error {
    client := common.GetRedisClient()
    key := driverLockPrefix + driverID
    return client.Del(common.Ctx, key).Err()
}

func reserveDriver(
    driverID string,
    orderID string,
) (bool, error) {
    client := common.GetRedisClient()
    key := driverLockPrefix + driverID
    ok, err := client.SetNX(
        common.Ctx,
        key,
        orderID,
        5*time.Minute,
    ).Result()
    return ok, err
}

func FindAndReserveNearestDriver(
    userLoc common.Location,
    orderID string,
) (string, error) {
    client := common.GetRedisClient()
    q := &redis.GeoRadiusQuery{
        Radius: 10000,
        Unit:   "m",
        Sort:   "ASC",
        Count:  10,
        WithDist: true,
    }
    locs, err := client.GeoRadius(
        common.Ctx,
        driversGeoKey,
        userLoc.Lng,
        userLoc.Lat,
        q,
    ).Result()
    if err == nil {
        for _, loc := range locs {
            ok, err := reserveDriver(loc.Name, orderID)
            if err != nil {
                continue
            }
            if ok {
                return loc.Name, nil
            }
        }
    }
    // Fallback: linear scan of all drivers
    drivers, derr := GetAllDriversFromStore()
    if derr != nil {
        return "", err
    }
    bestID := ""
    bestDsq := 0.0
    found := false
    for _, d := range drivers {
        free, _ := IsDriverAvailable(d.ID)
        if !free {
            continue
        }
        dx := d.Location.Lat - userLoc.Lat
        dy := d.Location.Lng - userLoc.Lng
        dsq := dx*dx + dy*dy
        if !found || dsq < bestDsq {
            bestDsq = dsq
            bestID = d.ID
            found = true
        }
    }
    if !found {
        return "", nil
    }
    ok, _ := reserveDriver(bestID, orderID)
    if ok {
        return bestID, nil
    }
    return "", nil
}

func IsDriverAvailable(
    driverID string,
) (bool, error) {
    client := common.GetRedisClient()
    key := driverLockPrefix + driverID
    n, err := client.Exists(
        common.Ctx, key,
    ).Result()
    if err != nil {
        return false, err
    }
    return n == 0, nil
}
