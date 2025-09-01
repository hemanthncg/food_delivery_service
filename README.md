# SwiftEats - Food Delivery Backend

SwiftEats is a modular monolith backend
for a food delivery service. It uses
Golang, Redis, and Kafka, and ships with
Docker Compose for a one-command local
run. A data simulator can generate live
driver location updates.

## Architecture (high-level)
- Modules: `apis`, `order`, `menu`,
  `driver`, `simulator`, `common`
- Storage: Redis (hashes/keys) for
  orders, drivers, restaurants, menus
- Streaming: Kafka topics for events:
  - `delivery_assignment`
  - `order_status_update`
- HTTP server: `:8080`
- See `ARCHITECTURE.md` for details

## Prerequisites
- Docker Desktop (or Docker Engine)
- Docker Compose plugin

## Quick start
1) Build and start stack
```bash
cd /Users/hemanthgajjarapu/food_delivery_service
docker compose up --build
```

This will start 4 containers:
- redis:7-alpine
- zookeeper:7.4.0
- kafka:7.4.0
- app (this project)

The app waits for Redis and Kafka to
accept connections before starting
Kafka consumers.

2) Verify health
```bash
curl -s http://localhost:8080/health
```
Expect: `OK`

## Environment configuration
The app reads core settings from envs
provided by Compose:

- `REDIS_ADDR` (default `redis:6379`)
- `KAFKA_BROKERS` (default `kafka:9092`)
- Simulator controls:
  - `SIMULATOR_ENABLED` (`1` or `0`)
  - `SIM_NUM_DRIVERS` (default `50`)
  - `SIM_EVENTS_PER_SEC` (default `10`)

Adjust in `docker-compose.yml` as needed.

## Demo: Place an order end-to-end
The steps below exercise restaurant
setup, driver setup, order placement,
and delivery status updates.

1) Enroll restaurant
```bash
curl -s -X POST \
  http://localhost:8080/restaurant/enroll \
  -H 'Content-Type: application/json' \
  -d '{
    "id": "rest1",
    "name": "Pizza Place",
    "location": {"lat": 19.11, "lng": 72.89},
    "menu": ["Margherita", "Pepperoni"]
  }'
```

2) Enroll a driver
```bash
curl -s -X POST \
  http://localhost:8080/driver/enroll \
  -H 'Content-Type: application/json' \
  -d '{
    "id": "driver1",
    "name": "Alice"
  }'
```

3) (Optional) Set driver location
```bash
curl -s -X POST \
  http://localhost:8080/driver/location \
  -H 'Content-Type: application/json' \
  -d '{
    "id": "driver1",
    "location": {"lat": 19.12, "lng": 72.88}
  }'
```

4) Place order
```bash
curl -s -X POST http://localhost:8080/order \
  -H 'Content-Type: application/json' \
  -d '{
    "id": "order1",
    "restaurant_id": "rest1",
    "items": ["Margherita"],
    "user_location": {"lat": 19.12, "lng": 72.88}
  }'
```
Expect response with `status: "created"`.
A background Kafka consumer will assign
nearest driver and publish a status.

5) Check order after a few seconds
```bash
curl -s \
  'http://localhost:8080/order/get?id=order1'
```
Expect to see `driver_id` set and
`status` similar to
`"Order assigned to delivery"`.

6) Mark delivered
```bash
curl -s -X POST \
  http://localhost:8080/order/delivered \
  -H 'Content-Type: application/json' \
  -d '{"order_id": "order1"}'
```
Then confirm:
```bash
curl -s \
  'http://localhost:8080/order/get?id=order1'
```
Expect `status: "Delivered"`.

## Postman collection
Import `POSTMAN_COLLECTION.json` into
Postman and set collection variable
`baseUrl` if needed (default is
`http://localhost:8080`).

## Useful endpoints
- Health: `GET /health`
- Restaurants: `POST /restaurant/enroll`
- Restaurants list: `GET /restaurants`
- Menu: `GET /menu?restaurant_id=...`
- Drivers: `POST /driver/enroll`
- Drivers list: `GET /drivers` (each item includes
  `available: true|false` derived from internal
  reservation lock)
- Driver location: `POST /driver/location`
- Get driver location: `GET /driver/location?id=...`
  (response includes `available`)
- Orders: `POST /order`
- Mark delivered: `POST /order/delivered`
- Get order: `GET /order/get?id=...`
- List orders: `GET /orders`

## Troubleshooting
- Kafka/Redis not ready:
  The app waits up to ~60s for both.
  If still failing, restart stack:
  ```bash
  docker compose down -v
  docker compose up --build
  ```
- Invalid request errors:
  Ensure JSON is valid and includes
  required fields. See Postman samples.
- Port conflicts:
  Stop other services using 8080/9092.

## Run tests locally (optional)
```bash
go test ./...
```
Integration-heavy tests are skipped.

## Clean up
```bash
docker compose down -v
```
Removes containers, networks, volumes.
