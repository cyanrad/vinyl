.PHONY: pb generate-db

generate-db:
	sqlc generate

reset-db:
	rm -f database.db
	make generate-db
	sqlite3 database.db < schema/schema.sql
