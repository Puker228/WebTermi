up:
	if docker network ls --format "{{.Name}}" | grep webtermi_network; then \
		echo "Network 'webtermi_network' already exists"; \
	else \
		echo "Creating network 'webtermi_network'..."; \
		docker network create webtermi_network; \
	fi

	docker compose build
	docker compose down
	docker compose up -d

stop:
	docker compose stop

down:
	docker compose down
	docker volume prune -a -f

prune:
	make stop
	docker compose down
	docker system prune -a -f
	docker volume prune -a -f

format:
	go fmt ./...

run:
	cd linux && docker build -t myub .
	go run cmd/app/main.go