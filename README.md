# Go E-Commerce API

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/Framework-Gin-008ECF?style=flat)](https://gin-gonic.com/)
[![Architecture](https://img.shields.io/badge/Architecture-Clean-success?style=flat)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

A robust, production-ready RESTful e-commerce API built with Go, adhering to **Domain-Driven Design (DDD)** and **Clean Architecture** principles. This project demonstrates a scalable structure for handling complex e-commerce workflows like user authentication, product management, and a dynamic shopping cart system.

## Architecture

The project follows the **Clean Architecture** (Hexagonal Architecture) pattern to ensure a high degree of decoupling, maintainability, and testability.

```mermaid
graph TD
    subgraph Presentation ["Presentation Layer (HTTP)"]
        Routes --> Controllers
    end

    subgraph Business ["Business Logic Layer"]
        Controllers --> UseCaseImpl["UseCase Implementation"]
    end

    subgraph Domain ["Domain Layer (Core)"]
        UseCaseImpl --> UseCasePort["UseCase Interface"]
        UseCaseImpl --> RepoPort["Repository Interface"]
        UseCaseImpl --> Entities
    end

    subgraph Infrastructure ["Infrastructure Layer"]
        RepoImpl["Repository Implementation"] -- "Implements" --> RepoPort
        RepoImpl --> GORM["GORM Models"]
        GORM --> DB[(PostgreSQL Database)]
    end

    style Domain fill:#f9f,stroke:#333,stroke-width:2px
    style Business fill:#bbf,stroke:#333,stroke-width:1px
    style Presentation fill:#dfd,stroke:#333,stroke-width:1px
    style Infrastructure fill:#eee,stroke:#333,stroke-width:1px
```

- **Entities**: Core business models and logic.
- **DTOs**: Data Transfer Objects for API requests and responses.
- **Ports (Interfaces)**: Abstractions for repositories and use cases.
- **Use Cases**: Orchestration of business rules and application logic.
- **Repositories**: Persistence logic using GORM (PostgreSQL).
- **Delivery (HTTP)**: API endpoints and controllers using the Gin framework.

## Key Features

### Authentication & Profile
- JWT-based authentication for secure access.
- Role-based access control (Customer vs. Seller).
- Multi-table registration (User + Customer/Seller profiles) using database transactions.
- Secure password hashing and token-based logout.
- Customer profile updates (Email, Phone, etc.).

### Category Management
- Hierarchical category management.
- Public read access and admin-protected (Seller/Admin) write operations.

### Product Catalog
- Full CRUD operations for products.
- Real-time stock tracking and inventory management.
- Public product browsing with detailed views.

### Shopping Cart System
- Persistent cart for authenticated customers.
- Real-time stock validation during "Add to Cart" and "Update Quantity".
- Support for batch item deletion.
- Auto-calculation of subtotals and total cart price.

## Technology Stack

- **Language**: [Go (Golang)](https://golang.org/)
- **Web Framework**: [Gin Gonic](https://gin-gonic.com/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Authentication**: JWT (JSON Web Tokens)
- **Testing**: [Testify](https://github.com/stretchr/testify) for mocking and assertions.
- **Migration**: SQL-based migrations for schema management.

## Getting Started

### 1. Prerequisites
- Go 1.25 or higher
- PostgreSQL 14+
- `golang-migrate` (optional, for running migrations via CLI)

### 2. Environment Configuration
Create a `.env` file in the root directory based on the project requirements:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=go_ecommerce
JWT_SECRET=your_super_secret_key
```

### 3. Database Migrations
The project uses SQL migrations located in `db/migrations/`. You can apply them using your preferred tool or the `migrate` CLI:

```bash
migrate -path db/migrations/ -database "postgres://user:pass@localhost:5432/db?sslmode=disable" up
```

### 4. Running the Application
```bash
# Install dependencies
go mod tidy

# Run the server
go run cmd/web/main.go
```

## API Endpoints

### Authentication
| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| POST | `/api/auth/register/customer` | Public | Register as a customer |
| POST | `/api/auth/register/seller` | Public | Register as a seller |
| POST | `/api/auth/login` | Public | Authenticate and get JWT |
| POST | `/api/auth/logout` | Protected | Invalidate current session |
| PUT | `/api/auth/customer` | Protected | Update customer profile |

### Product & Category
| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| GET | `/api/categories` | Public | List all categories |
| GET | `/api/categories/:id` | Public | Get category details |
| POST | `/api/categories` | Protected | Create new category |
| PUT | `/api/categories/:id` | Protected | Update category |
| DELETE | `/api/categories/:id` | Protected | Delete category |
| GET | `/api/products` | Public | List all products |
| GET | `/api/products/:id` | Public | Get product details |
| POST | `/api/products` | Seller | Add new product |
| PUT | `/api/products/:id` | Seller | Update product |
| DELETE | `/api/products/:id` | Seller | Delete product |

### Shopping Cart
| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| GET | `/api/cart` | Customer | View your shopping cart |
| POST | `/api/cart` | Customer | Add item to cart |
| PUT | `/api/cart/:id` | Customer | Update item quantity |
| DELETE | `/api/cart/batch` | Customer | Remove multiple items |

## Testing
The project includes comprehensive unit tests for business logic (Use Cases) and API controllers.

```bash
# Run all tests
go test ./...

# Run tests for a specific domain (e.g., Cart)
go test ./internal/usecase/cart/...
go test ./internal/delivery/http/cart/...
```

## Project Structure
```text
.
├── cmd/
│   └── web/            # Main entry point (main.go)
├── db/
│   └── migrations/     # PostgreSQL schema migrations
├── internal/
│   ├── delivery/       # HTTP Transport Layer (Gin Controllers & Routes)
│   ├── dto/            # Data Transfer Objects (Requests & Responses)
│   ├── entity/         # Domain Entities (Plain Go Structs)
│   ├── mocks/          # Mock Implementations for Unit Testing
│   ├── model/          # GORM Models (Persistence layer representation)
│   ├── port/           # Domain Interfaces (Repository & UseCase Ports)
│   ├── repository/     # Repository Implementations (PostgreSQL/GORM)
│   └── usecase/        # Business Logic / UseCase Implementation
├── .env.example        # Template for environment variables
└── README.md           # Documentation
```
