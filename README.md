# ⚔️ Army Builder API

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=for-the-badge&logo=docker)

A high-performance RESTful API and AI-ready backend for managing tabletop wargaming data. Built for Warhammer: Age of Sigmar (4th Edition) standards with a modular architecture ready for 40k.

## 🚀 Key Features

- **Relational Data Converter**: A sophisticated two-pass engine that transforms BattleScribe XML into clean, organized YAML.
- **Transactional Data Seeder**: A CLI tool that populates the database from YAML, utilizing ACID transactions and intelligent keyword de-duplication.
- **Recursive Stat Resolution**: Automatically resolves unit stats, weapons, and points across multiple linked library files.
- **Specialized Army Support**: First-class support for **Armies of Renown** (parent-linked) and **Regiments of Renown** (mercenaries).
- **Industry Standard Stats**: Supports complex stat strings (e.g., `5"`, `D3`, `3+`) to perfectly match official source material.
- **Deep Hydration**: API responses return fully nested unit data including Weapons, Abilities, Keywords, and Stat Modifiers.

## 🛠️ Tech Stack

- **Backend:** Go (Golang) 1.24+
- **Database:** PostgreSQL 16 (Relational, UUID-based)
- **Tooling:** SQLC (Type-Safe SQL), Docker Compose, Zap (Logging), YAML v3
- **Drivers:** pgx/v5 (Connection Pooling & Transactions)

## 🏗️ Project Structure

- `cmd/api/`: The main web server entry point.
- `cmd/converter/`: Two-pass XML to YAML transformation engine.
- `cmd/seeder/`: Transactional CLI tool for database ingestion.
- `internal/handlers/`: REST interface and JSON marshaling.
- `internal/services/`: Business logic and Army Validation engine.
- `internal/database/`: SQLC-generated type-safe database layer.
- `data/raw/`: Raw BattleScribe `.cat` and `.gst` source files.
- `data/factions/`: Organized YAML output, categorized by Game System and Army Type.

## 🔄 Data Pipeline

The project features a complete data pipeline to move from community-maintained XML to a high-performance relational database:

1. **Convert**: `cmd/converter` indexes all raw files to build a "Global Brain" of IDs, then performs a second pass to resolve links and output structured YAML.
2. **Organize**: Data is automatically sorted into `standard`, `armies_of_renown`, and `regiments_of_renown` subfolders.
3. **Seed**: `cmd/seeder` walks the organized directories and populates the PostgreSQL database, correctly linking parent/child faction relationships.

## 🚦 Getting Started

1. **Start Infrastructure**: `docker compose up -d`
2. **Convert Raw Data**: `go run ./cmd/converter` (Requires `.cat` files in `data/raw`)
3. **Seed the Database**: `go run ./cmd/seeder`
4. **Run the API**: `go run ./cmd/api`

## 🧪 Development & Testing

Run the full integration suite:
```bash
go test ./... -v
```

---
*Developed as a high-integrity backend foundation for future AI-driven list building.*
