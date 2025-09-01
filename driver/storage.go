package driver

import (
    "encoding/json"
    "food_delivery_service/common"
    "github.com/go-redis/redis/v8"
)

const (
    driverKey = "drivers"
    driversGeoKey = "drivers_geo"
)

func AddDriverToStore(d Driver) error {
    client := common.GetRedisClient()
    data, _ := json.Marshal(d)
    if err := client.HSet(
        common.Ctx,
        driverKey,
        d.ID,
        data,
    ).Err(); err != nil {
        return err
    }
    if d.Location.Lat != 0 || d.Location.Lng != 0 {
        _, _ = client.GeoAdd(
            common.Ctx,
            driversGeoKey,
            &redis.GeoLocation{
                Name: d.ID,
                Longitude: d.Location.Lng,
                Latitude: d.Location.Lat,
            },
        ).Result()
    }
    return nil
}

func GetAllDriversFromStore() ([]Driver, error) {
    client := common.GetRedisClient()
    vals, err := client.HVals(common.Ctx, driverKey).Result()
    if err != nil {
        return nil, err
    }
    var res []Driver
    for _, v := range vals {
        var d Driver
        if err := json.Unmarshal([]byte(v), &d); err == nil {
            res = append(res, d)
        }
    }
    return res, nil
}

func UpdateDriverLocationInStore(id string, loc common.Location) error {
    client := common.GetRedisClient()
    d, err := GetDriverFromStore(id)
    if err != nil {
        return err
    }
    d.Location = loc
    data, _ := json.Marshal(d)
    if err := client.HSet(
        common.Ctx,
        driverKey,
        id,
        data,
    ).Err(); err != nil {
        return err
    }
    _, _ = client.GeoAdd(
        common.Ctx,
        driversGeoKey,
        &redis.GeoLocation{
            Name: id,
            Longitude: loc.Lng,
            Latitude: loc.Lat,
        },
    ).Result()
    return nil
}

func GetDriverFromStore(id string) (Driver, error) {
    client := common.GetRedisClient()
    val, err := client.HGet(common.Ctx, driverKey, id).Result()
    if err != nil {
        return Driver{}, err
    }
    var d Driver
    if err := json.Unmarshal([]byte(val), &d); err != nil {
        return Driver{}, err
    }
    return d, nil
}
