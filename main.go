package main

import (
    "fmt"
    "log"
    "net/http"
    "food_delivery_service/apis"
    "food_delivery_service/order"
    "food_delivery_service/menu"
    "food_delivery_service/driver"
    "food_delivery_service/simulator"
    "food_delivery_service/common"
    "time"
    kafka "github.com/segmentio/kafka-go"
)

func main() {
    fmt.Println("Starting SwiftEats Backend...")
    cfg := common.LoadConfig()
    _ = cfg // placeholder for config usage
    apis.RegisterRoutes()
    // wait for dependencies inside container
    waitForRedis()
    waitForKafka()
    order.StartConsumers()
    order.Init()
    menu.Init()
    driver.Init()
    simulator.Init()
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func waitForRedis() {
    client := common.GetRedisClient()
    for i := 0; i < 60; i++ {
        if err := client.Ping(common.Ctx).Err(); err == nil {
            return
        }
        time.Sleep(1 * time.Second)
    }
}

func waitForKafka() {
    brokers := common.GetKafkaBrokers()
    addr := brokers[0]
    for i := 0; i < 60; i++ {
        conn, err := kafka.Dial("tcp", addr)
        if err == nil {
            _ = conn.Close()
            return
        }
        time.Sleep(1 * time.Second)
    }
}
