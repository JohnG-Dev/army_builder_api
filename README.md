# ⚔️ Army Builder API

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=for-the-badge&logo=docker)
![Status](https://img.shields.io/badge/Status-In%20Development-yellow?style=for-the-badge)

A robust RESTful API designed to manage tabletop wargaming data. This backend service allows users to query game systems, factions, units, and wargear, serving as the foundation for a future army list-building application.

Built with performance and maintainability in mind, utilizing **Clean Architecture** principles and **Type-Safe SQL**.

## 🏗️ Architecture & Design

This project follows a strict **Clean Architecture** pattern to ensure separation of concerns and testability:

*   **Handlers (`internal/handlers`)**: The HTTP transport layer. Responsible for parsing requests, validating inputs, and marshaling JSON responses.
*   **Services (`internal/services`)**: The core business logic layer. Orchestrates data flow and handles domain-specific rules.
*   **Database (`internal/database`)**: Type-safe Go code generated via **sqlc**. No ORM magic—just raw, performant SQL queries mapped to Go structs.
*   **State (`internal/state`)**: Dependency injection container for Database pools, Configuration, and structured Logging (Zap).

## 🛠️ Tech Stack

*   **Language:** Go (Golang)
*   **Database:** PostgreSQL
*   **Containerization:** Docker & Docker Compose
*   **SQL Generation:** [sqlc](https://sqlc.dev/) (Compiles SQL to type-safe Go)
*   **Routing:** Standard Library + Middleware
*   **Logging:** Uber Zap
*   **Migrations:** Golang-migrate

## 🗄️ Database Schema

The database is normalized to handle complex wargaming relationships:

*   **Games:** Top-level systems (e.g., Warhammer 40k, Age of Sigmar).
*   **Factions:** Sub-groups belonging to a game (e.g., Space Marines, Stormcast Eternals).
*   **Units:** The core data entity, containing stats (Move, Save, Wounds, Points).
*   **Relationships:**
    *   Units belong to Factions.
    *   Units have many **Weapons**.
    *   Units have many **Abilities** (Traits, Spells, Prayers).
    *   Units have many **Keywords** (Tags used for game logic).

## 🚀 Getting Started

### Prerequisites
*   Go 1.22+
*   Docker & Docker Compose

### Installation

1.  **Clone the repository**
    ```bash
    git clone https://github.com/john-1005/army_builder_api.git
    cd army_builder_api
    ```

2.  **Start the Database**
    Use Docker Compose to spin up the PostgreSQL container.
    ```bash
    docker-compose up -d
    ```

3.  **Run Migrations**
    Ensure your database schema is up to date.
    ```bash
    make migrate-up
    # OR manually with golang-migrate
    migrate -path sql/schema -database "postgres://postgres:postgres@localhost:5432/army_builder_api?sslmode=disable" up
    ```

4.  **Run the Server**
    ```bash
    go run .
    ```
    The server will start on port `:8080`.

## 🧪 Testing

The project emphasizes integration testing using a dedicated test database container.

*   **Integration Tests:** Validates the flow from Handler -> Service -> Database.
*   **Test Helpers:** Utilizes factory functions to minimize boilerplate in test files.

To run the full test suite:
```bash
go test ./... -v
```

## 🗺️ Roadmap & Status

### ✅ Phase 1: Core Infrastructure
- [x] Project structure and Clean Architecture setup.
- [x] Dockerized PostgreSQL environment.
- [x] Database Schema design and SQLC integration.

### ✅ Phase 2: API Implementation
- [x] CRUD endpoints for Games, Factions, Units, and Weapons.
- [x] Advanced filtering (filter units by keywords, points, etc.).
- [x] Structured error handling (RFC 7807 compliant).

### 🚧 Phase 3: Reliability (Current Focus)
- [ ] Achieve >80% code coverage on all Handlers.
- [ ] Implement comprehensive validation for edge cases.

### 🔜 Phase 4: Validation Engine
- [ ] Implement logic to validate army lists against points limits.
- [ ] Rules enforcement (min/max unit sizes, unique units).

### 🔜 Phase 5: Data Seeding
- [ ] Create YAML/JSON definition files for army lists.
- [ ] Build a seeder script to populate the DB with initial game data (e.g., Stormcast Eternals index).
