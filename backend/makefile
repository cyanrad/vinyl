.PHONY: pb generate-db

generate-db:
	sqlc generate

reset-db:
	rm -f ../files/database.db
	make generate-db
	sqlite3 ../files/database.db < schema/schema.sql
