.PHONY: pb generate-db

pb:
	pocketbase serve --dir ./pocketbase/database --publicDir ./pocketbase/public --migrationsDir ./pocketbase/migrations

generate-db:
	sqlc generate

reset-db:
	rm -f database.db
	make generate-db
	sqlite3 database.db < schema/schema.sql