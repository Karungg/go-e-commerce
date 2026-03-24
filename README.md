# Go E-Commerce API

A RESTful e-commerce API built with Go, Gin, GORM, and PostgreSQL following Hexagonal / Clean Architecture principles.

## Prerequisites

- Go 1.25+
- PostgreSQL
- Docker & Docker Compose (optional, for local development services)

## Environment Variables

Copy the example environment variables and configure them for your local setup:

```bash
cp .env.example .env
```
*(You will need to create a simple `.env` before running if you haven't yet, configuring things like `DB_HOST`, `DB_USER`, `JWT_SECRET`, etc.)*

## Running Locally

To run the application locally:

```bash
# Run the server
go run cmd/web/main.go
```

## Running with Docker

To spin up the required services (like the database) using Docker Compose:

```bash
docker-compose up -d
```
