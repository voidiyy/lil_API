http_server:
	$(MAKE) -C cmd/http_server run

MIGRATE_CMD = migrate -path internal/db_psql/database/migration/ -database "postgresql://typedb_psql:qwerzxcfgh2@localhost:5432/logic?sslmode=disable" -verbose

migrate-up:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down

migrate-force:
	$(MIGRATE_CMD) force 1


MARIA_MIGRATE_CMD = migrate -path internal/db_psql/database/migration/ -database "mysql://root:password@localhost:3306/MariaDB?sslmode=disable" -verbose

migrate-up:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down

migrate-force:
	$(MIGRATE_CMD) force 1
