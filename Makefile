
build-image:
	- docker build -t short_url .

run-app:
	- docker-compose -f docker-compose.yaml up

stop-app:
	- docker-compose -f docker-compose.yaml down