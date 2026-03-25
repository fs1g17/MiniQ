.PHONY: up
up:
	docker compose up -d --build
	sleep 3s
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=postgres password=password dbname=miniq host=0.0.0.0" goose -dir ./migrations up
