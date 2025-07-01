# SEA Catering 🥕 - Backend API

This repository contains the backend service for the **SEA Catering** application. It is a robust RESTful API built with Go (Golang) and the Fiber web framework, providing all the necessary business logic for user authentication, subscription management, and administrative reporting.

[![Language](https://img.shields.io/badge/Language-Go-blue.svg?style=for-the-badge&logo=go)](https://golang.org/)
[![Framework](https://img.shields.io/badge/Framework-Fiber-000000.svg?style=for-the-badge&logo=go)](https://gofiber.io/)
[![Database](https://img.shields.io/badge/Database-PostgreSQL-blue.svg?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
[![CI/CD](https://img.shields.io/badge/CI/CD-GitHub_Actions-2088FF.svg?style=for-the-badge&logo=github-actions)](https://github.com/features/actions)

---

## Table of Contents

- [About The Project](#about-the-project)
- [Features](#features)
- [API Documentation](#api-documentation)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Local Setup](#local-setup)
- [Makefile Commands](#makefile-commands)
- [CI/CD Pipeline](#cicd-pipeline)
- [Project Structure](#project-structure)
- [Contact](#contact)

## About The Project

This backend service is the powerhouse behind the SEA Catering platform. It handles all data persistence, user management, and core functionalities. It's designed to be scalable, secure, and efficient, exposing a clean REST API for web application to consume.

## Features

The API provides a complete set of features for a modern subscription-based service:

#### 🔐 Authentication & Authorization

- **User Registration:** Securely register new users.
- **JWT-Based Login:** Authenticate users and issue JSON Web Tokens (JWT) for session management.
- **Session Validation:** An endpoint to verify a user's token and retrieve their session data.
- **Role-Based Access Control (RBAC):** Differentiates between `USER` and `ADMIN` roles to protect sensitive endpoints.

#### 👨🏻 User-Facing Features

- **View Meal Plans:** Fetch a list of all available meal plans.
- **Create Subscriptions:** Subscribe to a meal plan with custom options (meal types, delivery days, allergies).
- **Manage Subscriptions:** View, update (e.g., pause/resume), and cancel personal subscriptions.
- **Submit Testimonials:** Provide feedback and ratings.

#### 👑 Admin-Facing Features

- **Subscription Reporting:** Generate business reports with date-range filters to view metrics.
- **Plan Management:** Update details of existing meal plans.

## API Documentation

This API is documented using Swagger 2.0.

- **Live Documentation:** You can access the live Swagger UI here: **[https://sea-catering-api.bccdev.id/docs](https://sea-catering-api.bccdev.id/docs)**

- **Local Documentation:** After running the server locally, the Swagger UI will be available at `http://localhost:3001/docs/index.html`.

To regenerate the documentation from code comments, you can use the following commands:

- **Using Make:**
  ```sh
  make generate-docs
  ```
- **Using Go CLI:**
  ```sh
  swag init -g cmd/api/main.go
  ```

## Tech Stack

- **Language:** [Go](https://golang.org/)
- **Web Framework:** [Fiber v2](https://gofiber.io/)
- **Database:** [PostgreSQL](https://www.postgresql.org/)
- **ORM:** [GORM](https://gorm.io/)
- **Containerization:** [Docker & Docker Compose](https://www.docker.com/)
- **Development Tool:** [Make](https://www.gnu.org/software/make/)
- **Deployment:** Ubuntu VM
- **CI/CD:** [GitHub Actions](https://github.com/features/actions)

## Getting Started

Follow these instructions to get the backend service running on your local machine for development.

### Prerequisites

- [Go](https://golang.org/doc/install) (v1.20 or later)
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)
- [Make](https://www.gnu.org/software/make/) (optional, but recommended for convenience)
- [swag](https://github.com/swaggo/swag) (for generating API docs)

### Local Setup

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/jevvonn/sea-catering-be.git
    cd sea-catering-be
    ```

2.  **Configure Environment Variables:**
    Create a `.env` file from the example. The default database values are pre-configured to work with the provided `docker-compose.yml`.

    ```bash
    cp .env.example .env
    ```

    Open the new `.env` file and fill all the values.

    ```env
    APP_ENV=development # or production
    APP_PORT=3001

    DB_USER=postgres
    DB_PASSWORD=yourpassword
    DB_NAME=sea_catering_db
    DB_HOST=localhost
    DB_PORT=5432

    # Generate a strong secret with: openssl rand -base64 32
    JWT_SECRET=your-super-strong-jwt-secret
    ```

3.  **Start the Database:**
    This command will start a PostgreSQL container in the background.

    ```bash
    docker-compose up -d
    ```

4.  **Install Go Dependencies:**

    ```bash
    go mod tidy
    ```

5.  **Initialize the Database (First-Time Setup Only):**
    For the first time, you need to run migrations to create the tables and then seed the database with initial data.

    **Step 5a: Run Migrations**

    - Using Make:
      ```sh
      make migrate-up
      ```
    - Using Go CLI:
      ```sh
      go run cmd/api/main.go -m up
      ```

    **Step 5b: Seed the Database**

    - Using Make:
      ```sh
      make db-seed
      ```
    - Using Go CLI:
      ```sh
      go run cmd/api/main.go -s
      ```

6.  **Run the Application:**
    Now you can start the API server.
    - Using Make:
      ```sh
      make run
      ```
    - Using Go CLI:
      ```sh
      go run cmd/api/main.go
      ```
      The server will now be running at `http://localhost:3001`.

## Makefile Commands

This project includes a `Makefile` to simplify common development tasks.

- `make run`: Starts the application server.
- `make migrate-up`: Applies all pending database migrations.
- `make migrate-down`: Drop all the database table and data.
- `make db-seed`: Seeds the database with initial data (e.g., user roles, meal plans).
- `make generate-docs`: Regenerates the Swagger API documentation.

## CI/CD Pipeline

This project uses a pull-based deployment strategy automated with GitHub Actions.

- **Workflow File:** `.github/workflows/deploy.yml`
- **Trigger:** The workflow is triggered on every `push` to the `main` branch.

#### Deployment Process:

The workflow connects to the production VM via SSH and executes a deployment script. The script performs the following steps on the server:

1.  Navigates to the project directory (`cd sea-catering-be`).
2.  Pulls the latest changes from the `main` branch (`git pull`).
3.  Uses Docker Compose to rebuild the application's Docker image and restart the service in the background (`sudo docker compose up -d --build`).

This approach ensures the application is always running the latest code from the `main` branch, with the build process happening directly on the target server.

#### Required GitHub Secrets:

For the deployment action to work, the following secrets must be configured in the GitHub repository settings:

- `SERVER_HOST`: The IP address or domain of the deployment server.
- `SERVER_USER`: The username for SSH login.
- `SSH_PRIVATE_KEY`: The private SSH key for authenticating with the server.
- `SERVER_PASSPHRASE`: The passphrase for the SSH key, if it has one.
- `SERVER_PORT`: The SSH port of the server (usually 22).

## Project Structure

The project follows a standard Go application layout to separate concerns.

```
/
├── cmd/
│   └── api/
│       └── main.go         # Main application entry point
├── config/                 # Environment variable and configuration loading
├── docs/                   # Swagger documentation files
├── internal/               # Main application logic
│   ├── controller/         # HTTP handlers (Fiber)
│   ├── dto/                # Data Transfer Objects for requests/responses
│   ├── domain/             # Domain models (entities and dtos)
│   ├── middleware/         # Custom Fiber middleware (e.g., auth)
│   ├── app/                # Application handler, usecases, and repositories
│   ├── service/            # Business logic
│   └── ...                 # Other internal packages
├── .github/workflows/      # CI/CD workflow definitions
├── docker-compose.yml      # Docker configuration for local database
└── go.mod                  # Go module dependencies
```

## Contact

Jevon Mozart - jmcb1602@gmail.com

Project Link: [https://github.com/jevvonn/sea-catering-be](https://github.com/jevvonn/sea-catering-be)
