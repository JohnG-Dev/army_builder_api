# ⚔️ Army Builder API

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=for-the-badge&logo=docker)

A high-performance RESTful API and AI-ready backend for managing tabletop wargaming data. Built for Warhammer: Age of Sigmar (4th Edition) standards.

## 🚀 Key Features

- **Transactional Data Seeder**: A CLI tool that populates the database from YAML, utilizing ACID transactions and intelligent keyword de-duplication.
- **Robust Validation**: An engine that enforces points limits, unit sizes, and keyword-specific army construction rules.
- **Industry Standard Stats**: Supports complex stat strings (e.g., `5"`, `D3`, `3+`) to perfectly match official source material.
- **Deep Hydration**: API responses return fully nested unit data including Weapons, Abilities, Keywords, and Stat Modifiers.

## 🛠️ Tech Stack

- **Backend:** Go (Golang) 1.24+
- **Database:** PostgreSQL 16 (Relational, UUID-based)
- **Tooling:** SQLC (Type-Safe SQL), Docker Compose, Zap (Logging), YAML v3
- **Drivers:** pgx/v5 (Connection Pooling & Transactions)

## 🏗️ Project Structure

- `cmd/api/`: The main web server entry point.
- `cmd/seeder/`: Transactional CLI tool for data ingestion.
- `internal/handlers/`: REST interface and JSON marshaling.
- `internal/services/`: Business logic and Army Validation engine.
- `internal/database/`: SQLC-generated type-safe database layer.
- `data/factions/`: YAML source files for game data.

## 🚦 Getting Started

1. **Start Infrastructure**: `docker-compose up -d`
2. **Seed the Library**: `go run ./cmd/seeder`
3. **Run the API**: `go run ./cmd/api`

## 🧪 Development & Testing

Run the full integration suite:
```bash
go test ./... -v
```

---
*Developed as a high-integrity backend foundation for future AI-driven list building.*
