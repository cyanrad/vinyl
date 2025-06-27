# created this cuz I keep forgetting it
generate-db:
	sqlc generate

# used for dev but should be removed or improved
reset-db:
	rm -f ../files/database.db
	make generate-db
	sqlite3 ../files/database.db < schema/schema.sql

run:
	go run . -media-path=../files
