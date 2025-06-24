run:
	rm -rf ./files/dist
	cd ./frontend && pnpm vite build --outDir ../files/dist/
	cd ./backend && go run .

build-win:
	rm -rf ./files/dist
	cd ./frontend && pnpm vite build --outDir ../files/dist/
	export GOOS=windows
	export GOARCH=amd64
	cd ./backend && go build -o ../vinyl.exe

build-lin:
# clear build files if exists 
	rm -rf ./files/dist
	rm -f ./files/database.db

# building frontend
	cd ./frontend && pnpm vite build --outDir ../files/dist/

# building backend 
	sqlite3 ./files/database.db < ./backend/schema/schema.sql
	cd ./backend && go build -o ../vinyl