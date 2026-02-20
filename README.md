```markdown
# Restaurant API

A robust, production-ready RESTful API built with **Go (Golang)** for managing restaurant operations. This project demonstrates scalable backend architecture, featuring secure authentication, role-based access control, and high-performance caching.

## 🚀 Features

* **High Performance**: Implements **Redis** for efficient caching and distributed rate limiting to handle high traffic loads.
* **User Management**: Secure registration (Customer & Seller), email activation, authentication, and profile management.
* **Role-Based Access Control (RBAC)**: Permission-based access for reading and writing data (e.g., `restaurant:read`, `restaurant:write`).
* **Restaurant Operations**: Create, update, delete, and list restaurants with advanced filtering and pagination.
* **Menu & Category System**: Organize food items into categories and menus linked to specific restaurants.
* **Security**: IP-based rate limiting (Token Bucket), graceful shutdowns, and secure password handling with bcrypt.
* **Mailing**: Integrated SMTP mailer for asynchronous user notifications.

## 🛠️ Tech Stack

* **Language**: Go (v1.23+)
* **Database**: PostgreSQL 15
* **Caching**: Redis 7
* **Router**: [httprouter](github.com/julienschmidt/httprouter)
* **Logging**: `log/slog` (Structured Logging)
* **Containerization**: Docker & Docker Compose
* **Configuration**: Command-line flags & Environment variables

## 🏗 Architecture & Performance

### ⚡ Redis Implementation
This project uses **Redis** to minimize database load and ensure scalability:
1.  **Authentication Caching**: User sessions and profiles are cached (`Cache-Aside` pattern). This avoids hitting PostgreSQL on every authenticated request, significantly reducing latency.
2.  **Distributed Rate Limiting**: Request counters are stored in Redis using a fixed-window algorithm. This allows the API to scale horizontally across multiple servers while maintaining accurate client limits.

## 📂 Project Structure

```text
.
├── cmd
│   └── api
│       ├── main.go         # Entry point, config, and dependency injection
│       ├── route.go        # HTTP route definitions
│       ├── server.go       # Server setup and graceful shutdown
│       ├── handlers.go     # HTTP handlers
│       ├── middleware.go   # Redis rate limiting, auth caching, and recovery
│       └── ...
├── internal
│   ├── models              # Database models and business logic
│   ├── mailer              # SMTP mailer implementation
│   └── validator           # Data validation helpers
└── go.mod                  # Module definition

```

## ⚡ Getting Started

### Prerequisites

* Go 1.23 or higher
* Docker & Docker Compose (Recommended)

### Configuration

The application is configured using command-line flags. You can also use environment variables to supply the DSN and Redis address (recommended for Docker).

| Flag | Env Variable | Default | Description |
| --- | --- | --- | --- |
| `-port` | `PORT` | `4000` | API server port |
| `-dsn` | `RESTAURANT_DB_DSN` | *(Required)* | PostgreSQL connection string |
| `-redis-addr` | `REDIS_ADDR` | `redis:6379` | Redis Host:Port |
| `-redis-password` | `REDIS_PASSWORD` | *(None)* | Redis Password |
| `-smtp-host` | `SMTP_HOST` | *(Required)* | SMTP host |
| `-limiter-enabled` | `LIMITER_ENABLED` | `true` | Enable rate limiter |
| `-limiter-rps` | `LIMITER_RPS` | `2` | Rate limiter requests per second |

## 🏃‍♂️ Running the Application

### Option 1: Using Docker Compose (Recommended)

This will spin up the API, PostgreSQL, and Redis containers simultaneously.

1. **Build and start the services:**
```bash
docker compose up -d --build

```


2. **Access the Application:**
* The API server will be available at `http://localhost:4000`.
* PostgreSQL is exposed on port `5435`.
* Redis is available internally to the app.


3. **Stop the services:**
```bash
docker compose down

```



### Option 2: Running Locally (Without Docker)

If you prefer running Go locally, ensure you have PostgreSQL and Redis instances running.

1. **Clone the repository and download dependencies:**
```bash
git clone [https://github.com/geekilx/restaurantAPI.git](https://github.com/geekilx/restaurantAPI.git)
cd restaurantAPI
go mod tidy

```


2. **Set Environment Variables:**
```bash
# Linux/macOS
export RESTAURANT_DB_DSN='postgres://user:pass@localhost/restaurant_api?sslmode=disable'
export REDIS_ADDR='localhost:6379'

# Run the app
go run ./cmd/api

```



## 🔗 API Endpoints

### Health Check

* `GET /v1/healthcheck` - Check API status and version.

### Users & Authentication

* `POST /v1/users` - Register a new customer.
* `POST /v1/users/login` - Authenticate and receive a token (Cached in Redis).
* `POST /v1/users/activate` - Activate a user account via token.
* `GET /v1/users/:id` - Get user details (Requires `restaurant:read`).

### Restaurants

* `GET /v1/restaurants` - List restaurants (supports pagination & filtering).
* `POST /v1/restaurants` - Create a new restaurant (Requires `restaurant:write`).
* `GET /v1/restaurants/:id` - Get a specific restaurant and its menu.

### Categories & Menus

* `GET /v1/category` - List all categories.
* `POST /v1/category` - Create a new category (Requires `restaurant:write`).
* `POST /v1/category/:id/menu` - Create a menu item under a category.

## 🤝 Contributing

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the branch (`git push origin feature/AmazingFeature`).
5. Open a Pull Request.
