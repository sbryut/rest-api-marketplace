# rest-api-marketplace

# Description
The project is a REST API service for a simple online marketplace, developed in the Go programming language. It provides basic functionality for managing users (registration, authorization) and advertisements (CRUD operations).

# Functionality
### Users
- Registration (Sigh Up): creating a new user account.
- Authorization (Sign In): logging into an existing account and receiving Access and Refresh tokens.
- Refresh Tokens: receiving a new pair of Access/Refresh tokens using an existing Refresh token.
### Advertisements
- Create Ad: adding a new advertisement by an authorized user.
- Get Ad By ID: viewing details of a specific advertisement.
- Get All Ads: viewing all advertisements with the ability to filter by price, sort by date/price and pagination.
- Update Ad: modify an existing ad by its owner.
- Delete Ad: delete an ad by its owner.

# Tech Stack
- Language: Go
- Web Framework: Echo
- Database: PostgreSQL
- DB Migrations: golang-migrate
- Password Hashing: bcrypt
- JWT Tokens: jwt-go
- Validation: go-playground/validator
- Logging: log/slog
- Linters: golangci-lint
- API Documentation: Swagger (swaggo/echo-swagger)

Create .env file in root directory and add following values:
```env
SERVER_HOST=localhost
SERVER_PORT=8080

POSTGRES_HOST=localhost
POSTGRES_USER=<your_user>
POSTGRES_PASSWORD=<your_password>
POSTGRES_EXTERNAL_PORT=5050
POSTGRES_DB=<your_db_name>

ENV_LOG=local
SIGNING_KEY=<random string>
ACCESS_TOKEN_TTL=3h
REFRESH_TOKEN_TTL=720h
```

