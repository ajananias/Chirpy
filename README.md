# Chirpy

Chirpy is lightweight HTTP server that I built in GO following the guide in the course "Learn HTTP servers in Go" from [boot.dev](http://boot.dev). It serves as the backend for a microblogging platform.

In this project I was able to implement authentication, login, authorizations, postings and database management sucessfully through HTTP requests.

## Why do I think it's a great project?

Chirpy demonstrates core backend development concepts including:

- RESTful API design
- HTTP request routing and middleware
- JSON serialization and deserialization
- Authentication and authorization
- Database persistence

It consider it to be a solid reference for me on how to build production-style HTTP servers in Go.

## Installation

### Prerequisites

- [Go](https://golang.org/dl/) 1.21 or later

### Clone the Repository

```bash
git clone https://github.com/ajananias/chirpy.git
cd Chirpy
```

### Install Dependencies

```bash
go mod tidy
```

### Running the server

```bash
go run .
```

The server should start on `http://localhost:8080` by default.

## API Endpoints

| Method | Endpoint            | Description                    |
|--------|---------------------|--------------------------------|
| GET    | /api/healthz        | Check server health            |
| GET    | /admin/metrics      | Check how many hits Chirpy has |
| POST   | /admin/reset        | Reset everything               |
| POST   | /api/users          | Create a new user              |
| POST   | /api/login          | Log in and get a token         |
| POST   | /api/chirps         | Create a new chirp             |
| GET    | /api/chirps         | Get all chirps                 |
| GET    | /api/chirps/{id}    | Get a chirp by ID              |
| DELETE | /api/chirps/{id}    | Delete a chirp                 |
| POST   | /api/refresh        | Create a new access token      |
| POST   | /api/revoke         | Revoke refresh token           |
| PUT    | /api/users          | Changes user membership status |
| POST   | /api/polka/webhooks | Validates membership upgrade   |
