# REFACTORED MYAPP
## Setup
### Createting ENV file

Create a .env file in the project root directory with the following variables (replace with your own values)

```bash
DB_DRIVER=
DB_NAME=
DB_HOST=
DB_PORT=
DB_USER=
DB_PASS=
SSL_MODE=

JWT_SECRET=

MONGO_DB =
MONGO_HOST=
MONGO_PORT=

TEST_DB_URL=
```

## Database Migrations
### Creating a Migration File
To create a new migration file, use the following command. This will create a new SQL migration file in the db/migrations directory:

```bash
goose -dir db/migrations create create_table_user sql
```

Replace create_table_user with a descriptive name for your migration.


### Running Migrations
To apply migrations to your PostgreSQL database, run:

```bash
goose -dir db/migrations postgres "postgres://user:password@host:port/dbname?sslmode=disable" up
```
Ensure you replace user, password, dbname, and host:port with your actual database credentials and connection details.

### Goose Help

For a complete list of Goose commands and options, you can run:

```bash
goose
```
This command will display all available Goose commands, making it easier to manage your database migrations.

## TESTING
### Running Unit Test Repository
```bash
go test ./test/unit/repository/... -coverpkg=./repository/... -coverprofile=coverage_repository.out
go tool cover -html=coverage_repository.out -o coverage_repository.html
```

### Running Unit Test Service
```bash
go test ./test/unit/service/... -coverpkg=./service/... -coverprofile=coverage_service.out
go tool cover -html=coverage_service.out -o coverage_service.html
```

### Running Integration Test Service
```bash
go test ./test/integrations/... -coverpkg=./handler/...,./service/...,./repository/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```