# HTTP Server with API to Database (MySQL/PostgreSQL)

This project implements an HTTP server with a REST API that can interact with either MySQL or PostgreSQL databases. The database type is configurable in the configuration file. The server supports JWT authentication, CSRF protection via cookies and headers, and includes validation for both.

## Features:
- Configurable database connection (supports MySQL and PostgreSQL).
- JWT authentication for secure API access.
- CSRF protection using cookies and headers.
- Configurable server settings.
- Ability to switch between MySQL and PostgreSQL databases.
- Use TLS cert&key for HTTPS.

Frontend and JS are very bad, thats just learning of golang API and DB, dont care about visual view

## Configuration Example:

```yaml
hostName: localhost
httpConf:
  listenAddr: "https://localhost:1234"
  certHostname: "localhost"
  timeOut: "4s"
  idleTimeout: "60s"
  certFile: "../../server.crt"
  keyFile: "../../server.key"
  jwtSecret: "your-jwt-secret"
dbConf:
  dbType: "mysql"  # Can be "postgres" for PostgreSQL
  dbURL: "your dbURL (for MySQL)"
```
Comands in Makefile for running server and DB migrations (only postgres;( )

```
http_server:
	$(MAKE) -C cmd/http_server run

MIGRATE_CMD = migrate -path internal/db_psql/database/migration/ -database "postgresql://typedb_psql:qwerzxcfgh2@localhost:5432/logic?sslmode=disable" -verbose

migrate-up:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down

migrate-force:
	$(MIGRATE_CMD) force 1
```
