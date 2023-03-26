
build-image:
	- docker build -t short_url .

up:
	- docker-compose up -d --build

down:
	- docker-compose down

cover:
	- go test -coverprofile=cover.out ./...
	- go tool cover -html=cover.out