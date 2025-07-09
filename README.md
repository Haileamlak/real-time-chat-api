## Real-Time Chat System in Go with Redis

A lightweight, scalable chat system built in Go, using **Redis** as the only data store. Supports:

* User sign-up and login
* Direct messaging (DM)
* Group chat
* Global broadcast messages
* Real-time communication via **WebSockets**
* Dockerized for easy local setup

---

## Architecture

This project follows the principles of **Clean Architecture**, promoting separation of concerns, testability, and scalability.

### Layers

* **`usecases/`** – Application-specific business rules
* **`infrastructure/`** – External services (e.g., Redis)
* **`repositories/`** – Data access layer (Redis)
* **`models/`** – Data structures and types
* **`controllers/`** – Web/HTTP/WS handlers
* **`routers/`** – HTTP routing and WebSocket handling
* **`main.go`** – Dependency injection & app startup

### Benefits

* Loose coupling between logic and frameworks
* Easy to test business logic independently
* Easy to swap out Redis, HTTP framework, or WebSocket lib

---

## Folder Structure

```
chat-system/
├── controllers/          # HTTP/WebSocket handlers
├── infrastructure/       # Redis connection and services
├── repositories/         # Data access layer
├── usecases/             # Business logic
├── models/               # Data models
├── main.go               # Application entry point
├── Dockerfile             # Docker configuration
├── docker-compose.yml     # Docker Compose setup
├── go.mod                # Go module dependencies
├── go.sum                # Go module checksums
├── README.md             # Project documentation
└── .env                  # Environment variables (optional)
```

---

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/haileamlak/chat-system.git
cd chat-system
```

### 2. Start with Docker Compose

```bash
docker-compose up --build
```

* Go server: `http://localhost:8080`
* Redis server: internal via `redis:6379`

> Requires [Docker Desktop](https://www.docker.com/products/docker-desktop/) with **WSL2 / Linux containers** enabled.

---

## WebSocket Testing

Connect to the WebSocket endpoint at:

```
ws://localhost:8080/ws?user=alice
```

* Use tools like [Postman](https://www.postman.com/), [Hoppscotch](https://hoppscotch.io/), or a browser client.
* Send messages in this format:

```json
{
  "type": "dm",         // "dm", "group", or "broadcast"
  "from": "alice",
  "to": "bob",          // or group name
  "content": "Hello!"
}
```

---

## Environment Variables

Set in `docker-compose.yml`:

```yaml
REDIS_ADDR=redis:6379
```

---

## Architecture Overview
    
* Real-time messaging is done via **WebSockets**
* All message routing uses **Redis Pub/Sub**
* Users and groups are stored in Redis hashes/sets
* Pub/Sub allows **horizontal scaling** with multiple Go instances

---

## TODO / Future Improvements

* Replace basic token auth with JWT
* Persist messages in Redis with TTL or logs
* Track online users (via Redis presence)
* Client-side reconnect logic
* Add unit/integration tests

---

## Author

Built by Haileamlak — feel free to contribute or fork!