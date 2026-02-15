# Restaurant API

A robust RESTful API built with Go (Golang) for managing restaurant operations, including user authentication, restaurant profiles, menus, and categories.

## 🚀 Features

* **User Management**: Registration (Customer & Seller), activation via email, authentication, and profile management.
* **Role-Based Access Control**: Permission-based access for reading and writing data (e.g., `restaurant:read`, `restaurant:write`).
* **Restaurant Management**: Create, update, delete, and list restaurants with filtering and pagination.
* **Menu & Category System**: Organize food items into categories and menus linked to specific restaurants.
* **Security**: Rate limiting, graceful shutdowns, and secure password handling with bcrypt.
* **Mailing**: Integrated SMTP mailer for user notifications (e.g., account activation).

## 🛠️ Tech Stack

* **Language**: Go (v1.25.4)
* **Database**: PostgreSQL
* **Router**: [httprouter](github.com/julienschmidt/httprouter)
* **Logging**: `log/slog` (Structured Logging)
* **Database Driver**: `lib/pq`
* **Configuration**: Command-line flags & Environment variables

## 📂 Project Structure

```text
.
├── cmd
│   └── api
│       ├── main.go         # Entry point, config, and dependency injection
│       ├── route.go        # HTTP route definitions
│       ├── server.go       # Server setup and graceful shutdown
│       ├── handlers.go     # HTTP handlers
│       ├── middleware.go   # Rate limiting, auth, and recovery middleware
│       └── ...
├── internal
│   ├── models              # Database models and business logic
│   ├── mailer              # SMTP mailer implementation
│   └── validator           # Data validation helpers
└── go.mod                  # Module definition

```

## ⚡ Getting Started

### Prerequisites

* Go 1.25 or higher
* PostgreSQL running locally or remotely

### Installation

1. **Clone the repository**
```bash
git clone [https://github.com/geekilx/restaurantAPI.git](https://github.com/geekilx/restaurantAPI.git)
cd restaurantAPI

```


2. **Download dependencies**
```bash
go mod download

```


3. **Database Setup**
Ensure your PostgreSQL database is running and create a database for the project.
```sql
CREATE DATABASE restaurant_api;
CREATE EXTENSION IF NOT EXISTS citext;

```


*(Note: Ensure you run any necessary migration SQL files to create the tables).*

### Configuration

The application is configured using command-line flags. You can also use environment variables to supply the DSN.

| Flag | Default | Description |
| --- | --- | --- |
| `-port` | `4000` | API server port |
| `-dsn` | `$RESTAURANT_DB_DSN` | PostgreSQL connection string |
| `-smtp-host` | `sandbox.smtp.mailtrap.io` | SMTP host |
| `-smtp-port` | `2525` | SMTP port |
| `-smtp-username` | `372d553c29c9c6` | SMTP username |
| `-smtp-password` | `3fb3fd1b008ee2` | SMTP password |
| `-smtp-sender` | `restaurantAPI ...` | Sender email address |
| `-limiter-enabled` | `true` | Enable rate limiter |
| `-limiter-rps` | `2` | Rate limiter requests per second |
| `-limiter-burst` | `4` | Rate limiter burst capacity |


## 🏃‍♂️ Running the Application

### Option 1: Using Docker Compose (Recommended)

You can easily spin up the application and its connected PostgreSQL database simultaneously using Docker. 

**Prerequisites:**
* Docker and Docker Compose installed on your system.

**Instructions:**
1. From the root of your project, build and start the containers:
   ```bash
   docker-compose up --build

```

2. The API server will be available at `http://localhost:4000`.
3. The PostgreSQL database will be mapped to port `5435` on your host machine.

To stop the containers and remove them, press `Ctrl+C` (if running attached) or run:

```bash
docker-compose down

```

### Option 2: Running Locally (Without Docker)

You can run the application directly using `go run`. It is recommended to set the DSN as an environment variable first.

**Linux/macOS:**

```bash
export RESTAURANT_DB_DSN='postgres://postgres:password@localhost/restaurant_api?sslmode=disable'
go run ./cmd/api

```

**Windows (PowerShell):**

```powershell
$env:RESTAURANT_DB_DSN='postgres://postgres:password@localhost/restaurant_api?sslmode=disable'
go run ./cmd/api

```
## 🔗 API Endpoints

### Health Check

* `GET /v1/healthcheck` - Check API status and version.

### Users & Authentication

* `POST /v1/users` - Register a new customer.
* `POST /v1/seller` - Register a new seller (grants `restaurant:write` permission).
* `POST /v1/users/activate` - Activate a user account via token.
* `POST /v1/users/authenticate` - Log in and retrieve an authentication token.
* `GET /v1/users/:id` - Get user details (Requires `restaurant:read`).
* `PATCH /v1/users/:id` - Update user details (Requires `restaurant:read`).
* `DELETE /v1/users/:id` - Delete a user (Requires `restaurant:read`).
* `PATCH /v1/resetpassword/:id` - Reset password (Requires `restaurant:read`).

### Restaurants

* `GET /v1/restaurants` - List restaurants (supports pagination & filtering).
* *Query Params:* `name`, `page`, `page_size`, `sort`


* `POST /v1/restaurants` - Create a new restaurant (Requires `restaurant:write`).
* *Body:* `{ "name": "String", "country": "String", "full_address": "String", "cuisine": "String", "status": "String" }`


* `GET /v1/restaurants/:id` - Get a specific restaurant and its menu.
* `PATCH /v1/restaurant/:id` - Update restaurant details (Requires `restaurant:write`).

### Categories & Menus

* `GET /v1/category` - List all categories.
* `POST /v1/category` - Create a new category (Requires `restaurant:write`).
* `GET /v1/category/:id` - Get all menus for a specific category.
* `GET /v1/restaurant/:id/categories` - Get categories for a specific restaurant.
* `GET /v1/menus` - List all menus.
* `POST /v1/category/:id/menu` - Create a menu item under a category (Requires `restaurant:write`).

## 🤝 Contributing

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the branch (`git push origin feature/AmazingFeature`).
5. Open a Pull Request.

## 📝 License

Distributed under the MIT License.

