package simulator

import (
    "math/rand"
    "os"
    "strconv"
    "strings"
    "time"
    "food_delivery_service/driver"
    "food_delivery_service/common"
)

// Init sets up the data simulator
func Init() {
    if !isEnabled() {
        return
    }
    go simulate()
}

// simulate generates N drivers and emits
// E location updates per second in total.
func simulate() {
    n := numDrivers()
    ids := genDriverIDs(n)
    for _, id := range ids {
        _ = driver.EnrollDriver(
            driver.Driver{
                ID:   id,
                Name: "Driver-" + id,
            },
        )
    }
    locs := make(map[string]common.Location)
    for _, id := range ids {
        locs[id] = common.RandomLocation()
    }
    rate := eventsPerSec()
    if rate < 1 {
        rate = 1
    }
    tick := time.NewTicker(
        time.Second / time.Duration(rate),
    )
    defer tick.Stop()
    rand.Seed(time.Now().UnixNano())
    for {
        <-tick.C
        id := ids[rand.Intn(len(ids))]
        loc := jitter(locs[id])
        locs[id] = loc
        _ = driver.UpdateLocation(id, loc)
    }
}

func jitter(l common.Location) common.Location {
    d := 0.0005
    l.Lat += (rand.Float64()*2 - 1) * d
    l.Lng += (rand.Float64()*2 - 1) * d
    return l
}

func isEnabled() bool {
    v := os.Getenv("SIMULATOR_ENABLED")
    if v == "" {
        return true
    }
    v = strings.ToLower(v)
    return v == "1" || v == "true"
}

func numDrivers() int {
    return envIntClamp(
        "SIM_NUM_DRIVERS", 50, 1, 50,
    )
}

func eventsPerSec() int {
    return envIntClamp(
        "SIM_EVENTS_PER_SEC", 10, 1, 1000,
    )
}

func envIntClamp(
    name string,
    def int,
    min int,
    max int,
) int {
    s := os.Getenv(name)
    if s == "" {
        return def
    }
    v, err := strconv.Atoi(s)
    if err != nil {
        return def
    }
    if v < min {
        return min
    }
    if v > max {
        return max
    }
    return v
}

func genDriverIDs(n int) []string {
    ids := make([]string, 0, n)
    for i := 1; i <= n; i++ {
        ids = append(ids, "driver"+itoa(i))
    }
    return ids
}

func itoa(i int) string {
    return strconv.Itoa(i)
}
