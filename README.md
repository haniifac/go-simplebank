# 🏦 Simple Bank with Golang

Welcome to the **Simple Bank**! This project is a simple implementation of a bank system that allows transferring money between accounts using **Golang**. It’s designed to be lightweight, efficient, and easy to extend. 🚀

---

## 💡 Features
- **Transfer Money** between accounts 🤑
- **Database Migration** with `golang-migrate`
- Auto-generated DB code using `SQLC`
- Built on the **Gin HTTP framework** for blazing-fast requests ⚡
- PostgreSQL-backed persistence
- Convenient `Makefile` to streamline development

---

## 📦 Dependencies
This project uses the following technologies:

1. [Golang Migrate](https://github.com/golang-migrate/migrate) - For database schema migrations.
2. [Gin](https://github.com/gin-gonic/gin) - A fast, easy-to-use HTTP framework.
3. [PostgreSQL](https://www.postgresql.org/) - Our reliable relational database of choice.
4. [SQLC](https://github.com/kyleconroy/sqlc) - Generates type-safe Go code from SQL queries.
5. [Makefile](https://www.gnu.org/software/make/manual/make.html) - Automates common tasks so you can focus on coding.

---

## 🚀 Getting Started
Follow these steps to get up and running with the Simple Bank.

### 1. Clone the Repo
```bash
git clone https://github.com/your-username/simple-bank.git
cd simple-bank
```

### 2. Get and verify the golang dependencies
```bash
go mod tidy
go mod verify
```

Install makefile and go-migrate. This step is crucial to perform the `make` command

### 3. Start PostgreSQL with Docker
Run the following Makefile command to spin up a PostgreSQL container:
```bash
make postgres
```

### 4. Create the database
```bash
make createdb
```

### 5. Apply the migration schema
```bash
make migrateup
```

### 6. Run the app
```bash
go run main.go
```

## Makefile usage examples

| Command                                               | Description                                                       |
| :---------------------------------------------------: | :---------------------------------------------------------------: |
| `make postgres`                                       | Start PostgreSQL database using Docker                            |
| `make createdb`                                       | Create a new database inside the running PostgreSQL container     |
| `make dropdb`                                         | Drop (delete) the database                                        |
| `make new_migrate name=![#c5f015]your_migration_name` | Create a new migration file                                       |
| `make migrateup`                                      | Apply all new database migrations                                 |
| `make migrateup count=![#c5f015]N`                    | Apply N migrations                                                |
| `make migratedown`                                    | Rollback all migrations                                           |
| `make migratedown count=![#c5f015]N`                  | Rollback N migrations                                             |
| `make sqlc`                                           | Generate Go code from SQL queries using SQLC                      |
| `make test`                                           | Run unit tests with coverage                                      |
| `make server`                                         | Run the Go server (alternative to `go run main.go`)               |
| `make mock`                                           | Generate mocks for unit tests using `mockgen`                     |
