## Project Structure

This document explains the repository
layout, the purpose of each folder, and
key modules used by the system.

### Repository layout

- `main.go`: app entrypoint, wires routes
  and starts background consumers.
- `ARCHITECTURE.md`: high level
  architecture overview and rationale.
- `README.md`: how to build and run
  the system locally.
- `POSTMAN_COLLECTION.json`: Postman
  requests for all API endpoints.
- `docker-compose.yml`: local infra for
  Redis, Zookeeper, Kafka, and the app.
- `Dockerfile`: multi stage build for the
  application binary.
- `go.mod`, `go.sum`: Go module files.

### Folders

#### `apis/`
- Purpose: HTTP API layer. Defines all
  handlers and registers routes.
- Key files:
  - `handlers.go`: exported HTTP handlers.
    Includes `HealthHandler`, restaurant,
    menu, driver, and order endpoints.
  - `routes.go`: binds URL paths to
    handlers in `handlers.go`.
  - `api_test.go`: unit tests for
    selected handlers (minimal; some
    integration tests are skipped).

#### `order/`
- Purpose: Order lifecycle, assignment,
  and persistence.
- Key files:
  - `order.go`: `Order` model and `Init`.
  - `service.go`: business logic for
    placing orders, driver assignment,
    listing and fetching orders.
  - `storage.go`: Redis backed store for
    orders (hash `orders`).
  - `consumers.go`: Kafka consumers for
    `delivery_assignment` and
    `order_status_update` topics.
  - `order_test.go`: legacy test stub.

#### `menu/`
- Purpose: Restaurant enrollment and menu
  storage and retrieval.
- Key files:
  - `menu.go`: core types `Restaurant`,
    `Menu` and in memory samples.
  - `service.go`: higher level APIs for
    enrolling and listing restaurants,
    and for fetching menus.
  - `storage.go`: Redis backed store for
    restaurants (hash `restaurants`) and
    menus (keys `menu:<restaurant_id>`).
  - `menu_test.go`: minimal test stub.

#### `driver/`
- Purpose: Driver enrollment, location
  tracking, and nearest driver lookup.
- Key files:
  - `driver.go`: `Driver` model and `Init`.
  - `service.go`: enroll, list, get and
    update location, closest driver. Also
    maintains availability via a Redis
    reservation lock (`driver_lock:<id>`),
    exposing `IsDriverAvailable` and
    atomic reservation APIs used by order
    assignment.
  - `storage.go`: Redis backed store for
    drivers (hash `drivers`) and GEO index
    (`drivers_geo`) for nearest lookups.
  - `driver_test.go`: minimal test stub.

#### `simulator/`
- Purpose: Generate driver data for local
  testing. Can enroll up to 50 drivers
  and emit ~10 location updates per sec
  across all drivers (defaults).
- Key files:
  - `simulator.go`: reads env settings,
    enrolls drivers, and sends location
    updates with a random walk.
  - Env variables:
    - `SIMULATOR_ENABLED` (default true)
    - `SIM_NUM_DRIVERS` (default 50)
    - `SIM_EVENTS_PER_SEC` (default 10)

#### `common/`
- Purpose: Shared utilities for config,
  Redis, Kafka, and helper types.
- Key files:
  - `common.go`: `Config` and loader.
  - `redis.go`: singleton Redis client.
    Reads `REDIS_ADDR` (defaults to
    `redis:6379`).
  - `kafka.go`: Kafka helpers and
    consumer loop. Reads `KAFKA_BROKERS`
    (defaults to `kafka:9092`).
  - `gps.go`: `Location` type and a
    random location generator.

### Runtime wiring

- `main.go`:
  - Registers routes from `apis/`.
  - Waits for Redis and Kafka to be
    reachable before starting consumers.
  - Starts order consumers and the
    simulator goroutine.
  - Serves HTTP on `:8080`.

### Local infrastructure

- `docker-compose.yml` provisions:
  - Redis on `redis:6379`.
  - Zookeeper for Kafka.
  - Kafka broker on `kafka:9092`.
  - The app with env wiring:
    `REDIS_ADDR`, `KAFKA_BROKERS`, and
    simulator envs.

### Notes

- Modules are decoupled with clean
  interfaces. Most cross module calls
  are function calls, not RPC.
- Kafka is used to decouple order
  placement and driver assignment.
