
build-image:
	- docker build -t short_url .

run-app:
	- docker-compose -f docker-compose.yaml up

stop-app:
	- docker-compose -f docker-compose.yaml down

cover:
	- go test -coverprofile=cover.out ./...
	- go tool cover -html=cover.out