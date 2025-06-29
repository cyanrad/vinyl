generate-db:
	cd ./backend && sqlc generate

# clear build files if exists 
clear-build:
	rm -rf ./files/dist
	rm -f ./vinyl
	rm -f ./files/database.db
	rm -f ./vinyl.exe

build-frontend:
	cd ./frontend && pnpm vite build --outDir ../files/dist/

build-db:
	sqlite3 ./files/database.db < ./backend/schema/schema.sql

build-all:
	$(MAKE) clear-build
	$(MAKE) build-frontend
	$(MAKE) build-db

build-backend-linux:
	cd ./backend && go build -o ../vinyl

build-backend-windows:
	cd ./backend && \
		GOOS="windows" \
		GOARCH="amd64" \
		CGO_ENABLED="1"\
		CC="x86_64-w64-mingw32-gcc" \
		CXX="x86_64-w64-mingw32-g++" \
		go build -o ../vinyl.exe

clean-run:
	$(MAKE) build-all
	cd ./backend && go run . -media-path=../files -ingest -source=local
	cd ./backend && go run . -media-path=../files

run:
	cd ./backend && go run . -media-path=../files

build-win-on-linux:
	$(MAKE) build-all
# we can't use the output executable to ingest
	$(MAKE) build-backend-windows
	cd ./backend && go run . -ingest -source=local

build-linux:
	$(MAKE) build-all
	$(MAKE) build-backend-linux
	./vinyl -ingest

# NOTE: this is used for dev and should not be used for prod
bundle:
	zip -r release.zip vinyl.exe files/
