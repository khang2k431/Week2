#  Task Manager API (Gin + GORM + PostgreSQL)

## Introduction
This is a REST API built with **Go**, using **Gin Framework** and **GORM** for a simple **Task Management System**.  

### Main Features
- User registration & login with **JWT Authentication**  
- CRUD operations for **Tasks** (create, read, update, delete)  
- Role-based access control (**Admin / User**)  
- Rate limiting to prevent abuse  
- User action logging  

---

##  Project Structure
Week2/
├── config/ # Database & env config
│ └── config.go
├── controllers/ # Request/response handlers
│ ├── auth.go
│ └── task.go
├── middlewares/ # JWT, Rate limiting
│ ├── auth_middleware.go
│ └── rate_limit.go
├── models/ # User, Task models
│ └── models.go
├── utils/ # JWT utility functions
│ └── jwt.go
├── main.go # Entry point
├── go.mod
└── README.md

---

##  Setup

### 1. Clone the project
```bash
git clone https://github.com/yourusername/taskmanager.git
cd taskmanager/Week2

Create .env file
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=taskdb
DB_PORT=5432
SERVER_PORT=8080
JWT_SECRET=supersecret

Install dependencies
go mod tidy

Run the server
go run main.go

