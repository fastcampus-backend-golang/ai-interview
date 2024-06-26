docker:
	docker compose up -d
	
build:
	go build -o interview .

run: build
	./interview