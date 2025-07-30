# 🧾 Ledger Service

A containerized banking ledger service built with Golang. This project uses PostgreSQL for account storage, MongoDB for ledger entries, and RabbitMQ for asynchronous transaction processing.

---

## 🧱 Tech Stack

- **Go 1.22**
- **PostgreSQL** (Account storage)
- **MongoDB** (Ledger entries)
- **RabbitMQ** (Async transaction processing)
- **Docker & Docker Compose**

---

## 🚀 Getting Started

### ✅ Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

---

## 📦 Running the Application

### 1. Clone the Repository

```bash
git clone https://github.com/berbreik/Ledger.git
cd ledger
```

### 2. Build and Run with Docker Compose

```bash
docker-compose up --build
```
### 3. Access the Services
- **PostgreSQL**: `localhost:5432`
- **MongoDB**: `localhost:27017`
- **RabbitMQ**: `localhost:15672` (default credentials: `guest/guest`)
- **Ledger Service**: `localhost:8080`

