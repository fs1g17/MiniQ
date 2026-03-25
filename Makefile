.PHONY: up
up:
	docker compose up -d --build
	sleep 3s
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres password=password dbname=miniq host=localhost port=5432 sslmode=disable" goose -dir ./migrations up
