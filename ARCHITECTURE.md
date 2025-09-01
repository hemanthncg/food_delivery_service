# ARCHITECTURE.md

## Overview
This document describes the modular
monolith architecture for the SwiftEats food delivery
backend. The system is designed for local deployment
using Docker, with a focus on scalability, resilience,
and performance. The architecture leverages Golang,
SQL, Redis, Apache Kafka, and in-memory caching.

## System Modules
The backend is organized into the following modules:

### 1. API Gateway
- Handles all HTTP/gRPC requests.
- Routes requests to internal modules.
- Performs basic request validation.

### 2. Order Management
- Manages order lifecycle (create, update, track).
- Handles order state transitions.
- Interacts with mocked payment logic.

### 3. Menu & Restaurant
- Manages restaurant and menu data.
- Provides fast menu/status queries.
- Uses Redis and in-memory cache for speed.

### 4. Driver & Logistics
- Tracks driver status and location.
- Ingests GPS data from drivers.
- Publishes/consumes location events via Kafka.

### 5. Analytics
- Consumes real-time events from Kafka.
- Aggregates data for reporting and monitoring.

### 6. Data Simulator
- Generates mock driver GPS data.
- Simulates up to 50 drivers for local testing.

## Technology Choices
- **Golang:** Main backend language.
- **SQL Database:** Stores core business data.
- **Redis:** Caching for menus, sessions, etc.
- **In-Memory Cache:** For ultra-fast, ephemeral data.
- **Apache Kafka:** Event streaming for driver data.
- **Docker:** Containerizes all components.

## Technology Justification

### Monolith
- Enables faster initial development.
- Easier to deploy and debug locally.
- Reduces operational overhead for small teams.

### Golang
- High performance for backend workloads.
- Excellent for handling data concurrency.
- Well-supported and efficient for microservices or monoliths.

### Kafka
- Simple and robust event streaming platform.
- Ideal for real-time data pipelines and decoupling modules.

### In-Memory Caching
- Boosts performance by reducing DB load.
- Enables fast data access for user-facing features.
- Improves system responsiveness under load.

## Data Flow
1. **Order Placement:**
   - API Gateway → Order Management
   - Order Management updates DB, mocks payment
2. **Menu Fetch:**
   - API Gateway → Menu Module
   - Menu Module checks Redis, then DB
3. **Driver Location:**
   - Driver Simulator → Kafka
   - Driver Module consumes Kafka, updates Redis
   - API Gateway fetches live location from Redis
4. **Analytics:**
   - Analytics Module consumes Kafka events
   - Aggregates and stores analytics data

## Rationale
- **Modular Monolith:**
  - Simple to develop and deploy locally.
  - Clear module boundaries for future scaling.
- **Redis & In-Memory Cache:**
  - Meet strict response time targets for menu/status.
- **Kafka:**
  - Decouples real-time data ingestion and analytics.
- **Docker:**
  - Ensures consistent local development.

## Scalability & Maintainability
- Each module is a package with clear interfaces.
- Modules communicate via function calls, not network.
- Easy to split modules into microservices later.

## Fault Tolerance
- Non-critical module failures do not affect core flow.
- Kafka and Redis are run as local containers.
- Payment and restaurant integrations are mocked.

## Testing
- Data simulator validates real-time flows.
- Unit and integration tests for all modules.

---
This architecture meets the current business and
technical requirements, and is designed for easy
future evolution.
