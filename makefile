# clear build files if exists 
clear-build:
	rm -rf ./files/dist
	rm -f ./files/database.db

build-frontend:
	cd ./frontend && pnpm vite build --outDir ../files/dist/

build-db:
	sqlite3 ./files/database.db < ./backend/schema/schema.sql

build-all:
	$(MAKE) clear-build
	$(MAKE) build-frontend
	$(MAKE) build-db

clean-run:
	$(MAKE) build-all
	cd ./backend && go run .

run:
	cd ./backend && go run .

build-win:
	$(MAKE) build-all
	export GOOS=windows
	export GOARCH=amd64
	cd ./backend && go build -o ../vinyl.exe

build-lin:
	$(MAKE) build-all
	cd ./backend && go build -o ../vinyl

