# ✨ Personal Website Backend

A modern, secure REST API built with Go to power your personal website, blog, and portfolio.

![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Fiber](https://img.shields.io/badge/Fiber-Web_Framework-00ACD7?style=for-the-badge)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?style=for-the-badge&logo=postgresql&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)

## 📋 Overview

This backend API provides a robust foundation for your personal website with blog posts, portfolio items, and secure user authentication. Built with clean architecture principles and modern Go practices.

## 🚀 Features

- 🔐 **Secure Authentication** — JWT-based auth with Argon2id password hashing
- 🛡️ **Admin Dashboard** — Protected routes for content management
- 📝 **Blog Management** — Create, update, and publish articles
- 🖼️ **Portfolio System** — Showcase your work with detailed portfolio items
- 📦 **API Versioning** — Future-proof your API with versioned endpoints
- 🔒 **Security First** — CORS protection, rate limiting, secure headers
- 📊 **Structured Logging** — Comprehensive logging with Zap
- 🧩 **Clean Architecture** — Separation of concerns for maintainability

## 🛠️ Technology Ecosystem

### 🔩 Core Infrastructure

- **🔷 Go 1.24+** — High-performance, statically typed language for scalable backends
- **🐘 PostgreSQL** — Enterprise-grade relational database with advanced features
- **🌐 RESTful API** — Industry-standard architecture for web services
- **🧪 Unit Testing** — Comprehensive test coverage for reliability

### 📚 Powerful Libraries

- **⚡ Fiber** — Lightning-fast, Express-inspired web framework optimized for performance
- **🦢 Goose** — Elegant database migration management system
- **🧬 Native SQL** — Hand-crafted queries for precise database operations
- **🔑 JWT** — Secure token-based authentication protocol
- **⚡ Zap** — Blazingly fast, structured logging framework
- **🔄 Argon2id** — State-of-the-art cryptographic hashing algorithm
- **🛡️ CORS** — Cross-Origin Resource Sharing protection
- **⏱️ Rate Limiter** — Traffic control mechanism for API stability

## 📂 Project Structure

```
├── cmd/               # Application entry points
│   ├── api/           # Main API server
│   └── hash/          # Password hashing utility
├── config/            # Application configuration
├── db/                # Database connection and migrations
├── internal/          # Internal packages
│   ├── controller/    # HTTP request handlers
│   ├── logger/        # Logging infrastructure
│   ├── middleware/    # HTTP middleware components
│   ├── model/         # Data models and DTOs
│   ├── repository/    # Data access layer
│   ├── router/        # Route definitions
│   ├── service/       # Business logic layer
│   └── util/          # Utility functions
└── scripts/           # Helper scripts
```

## 🌐 API Endpoints

### 🔓 Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/public/articles` | List published articles |
| `GET` | `/api/v1/public/articles/:id` | Get article by ID |
| `GET` | `/api/v1/public/articles/slug/:slug` | Get article by slug |
| `GET` | `/api/v1/public/portfolios` | List published portfolios |
| `GET` | `/api/v1/public/portfolios/:id` | Get portfolio by ID |
| `GET` | `/api/v1/public/portfolios/slug/:slug` | Get portfolio by slug |

### 🔑 Auth Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/auth/login` | Login and receive JWT tokens |

### 🔒 Admin Endpoints (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/admin/profile` | Get user profile |
| `PUT` | `/api/v1/admin/profile` | Update user profile |
| `PUT` | `/api/v1/admin/profile/avatar` | Update profile avatar |
| `PUT` | `/api/v1/admin/profile/password` | Change password |
| `GET` | `/api/v1/admin/articles` | List all articles (including drafts) |
| `POST` | `/api/v1/admin/articles` | Create new article |
| `PUT` | `/api/v1/admin/articles/:id` | Update existing article |
| `DELETE` | `/api/v1/admin/articles/:id` | Delete article |
| `GET` | `/api/v1/admin/portfolios` | List all portfolios (including drafts) |
| `POST` | `/api/v1/admin/portfolios` | Create new portfolio |
| `PUT` | `/api/v1/admin/portfolios/:id` | Update existing portfolio |
| `DELETE` | `/api/v1/admin/portfolios/:id` | Delete portfolio |

## 🏁 Getting Started

### Prerequisites

- 🔷 Go 1.24 or later
- 🐘 PostgreSQL 14 or later

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/personal-website-backend.git
   cd personal-website-backend
   ```

2. **Set up your environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Create your PostgreSQL database**
   ```bash
   createdb personal_website
   ```

4. **Install dependencies**
   ```bash
   go mod download
   ```

5. **Run the application**
   ```bash
   go run cmd/api/main.go
   ```

The server will start on the port specified in your `.env` file (default: 8080).

### 🗄️ Database Migrations

Migrations are automatically run when the application starts, but you can also:

```bash
# Create a new migration
go run cmd/api/main.go db:create name_of_migration

# Apply all pending migrations
go run cmd/api/main.go db:migrate

# Roll back the most recent migration
go run cmd/api/main.go db:rollback

# Reset the entire database
go run cmd/api/main.go db:reset
```

### 🔐 Default Admin User

The system comes with a default admin user:
- **Username**: `admin`
- **Password**: `admin`

⚠️ **Important**: Change this password immediately in a production environment!

You can generate a new password hash using:
```bash
go run cmd/hash/main.go your_secure_password
```

## 🛡️ Security Features

- 🔒 **JWT Authentication** — Secure token-based auth with refresh tokens
- 🔑 **Argon2id Hashing** — Modern, secure password hashing
- 🛡️ **CORS Protection** — Configurable cross-origin resource sharing
- ⏱️ **Rate Limiting** — Protect against brute-force and DDoS attacks
- 🔒 **Secure Headers** — HTTP security headers (HSTS, CSP, etc.)
- 🔍 **Input Validation** — Request validation to prevent injection attacks
- 📊 **Structured Logging** — Comprehensive logging with sensitive data redaction

## 📜 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

Built with ❤️ using Go and modern web technologies.
