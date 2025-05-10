# FlatSync API

This is a backend API for the **FlatSync** application built using **Go** and the **Fiber** framework. It provides user authentication (register/login), JWT-based token authentication, and access to protected routes.

## Prerequisites

To run this project locally, you'll need the following:

* **Go** (1.18+)
* **SQLite3** (for local database)
* **Docker** (optional, for containerization)
* **Postman** or **curl** (for testing API endpoints)

## Setup Instructions

Follow these steps to set up and run the project locally.

### 1. Clone the repository

First, clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/FlatSync-API.git
cd FlatSync-API
```

### 2. Install Go dependencies

Install all the Go dependencies required for the project:

```bash
go mod tidy
```

This will automatically download and install the required packages and dependencies.

### 3. Set up environment variables

Create a `.env` file in the root of your project directory. You can copy the template from `.env.example` or create it manually. This file contains sensitive information such as your JWT secret key.

Example `.env` file:

```
JWT_SECRET=your_jwt_secret_key_here
```

> **Note:** Replace `your_jwt_secret_key_here` with a secure secret key for signing JWT tokens.

### 4. Set up the database

This project uses **SQLite** as the database. To set up the database:

1. You can use any SQLite management tool to view and interact with the `db.sqlite3` file.
2. If you need to create the database manually, you can run the following SQL commands:

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    email TEXT UNIQUE,
    password TEXT,
    dob TEXT,
    avatar TEXT
);
```

Alternatively, the app will automatically create the `users` table when it starts up, if it doesn't exist yet.

### 5. Run the server

After setting up the environment and database, you can run the server with:

```bash
go run main.go
```

This will start the server on `http://localhost:3000` by default.

### 6. Test the endpoints

#### Register a new user

To register a new user, send a **POST** request to `/register` with the following JSON body:

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword",
  "dob": "1990-01-01",
  "avatar": "https://example.com/avatar.png"
}
```

#### Login a user

To log in a user, send a **POST** request to `/login` with the following JSON body:

```json
{
  "email": "john@example.com",
  "password": "securepassword"
}
```

If the credentials are valid, the server will return a JWT token.

#### Access a protected route

To access a protected route, send a **GET** request to `/protected` with the **Authorization** header containing the JWT token:

```
Authorization: Bearer <your_jwt_token_here>
```

### 7. Testing with Postman or Curl

You can use Postman or curl to test the API endpoints.

For example, using curl to register a user:

```bash
curl -X POST http://localhost:3000/register \
    -H "Content-Type: application/json" \
    -d '{"name": "John Doe", "email": "john@example.com", "password": "securepassword", "dob": "1990-01-01", "avatar": "https://example.com/avatar.png"}'
```

### 8. Error handling

If an error occurs, the server will return an appropriate status code (e.g., `400`, `401`, `500`) and a JSON response indicating the error.

---

## Additional Notes

* **JWT Secret**: The JWT secret is important for signing and verifying the tokens. It should be kept secure and should not be shared publicly.

* **Database**: The database used in this project is SQLite for simplicity. If you're planning to use this app in production, you may need to migrate to a more robust database like MySQL or PostgreSQL.
