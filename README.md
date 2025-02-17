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
5. `Makefile` - Automates common tasks so you can focus on coding.

---

## 🚀 Getting Started
Follow these steps to get up and running with the Simple Bank.

### 1. Clone the Repo
```bash
git clone https://github.com/your-username/simple-bank.git
cd simple-bank
```

### 2. Start PostgreSQL with Docker
Run the following Makefile command to spin up a PostgreSQL container:
```bash
make postgres
```
