
build-image:
	- docker build -t short_url .

run-app:
	- docker-compose -f docker-compose-app.yaml up

stop-app:
	- docker-compose -f docker-compose-app.yaml down

run-env:
	- docker-compose -f docker-compose-env.yaml up

stop-env:
	- docker-compose -f docker-compose-env.yaml down


cover:
	- go test -coverprofile=cover.out ./...
	- go tool cover -html=cover.out