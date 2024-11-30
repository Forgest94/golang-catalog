.EXPORT_ALL_VARIABLES:
PROJECT_NAME := catalog
NET_IP := 10.227.100

install:
	mkdir -p _db/elasticsearch
	cp .env.example .env
	sudo chown -R 1000 _db/elasticsearch
	sudo sh -c "grep -qxF '${NET_IP}.107    elasticsearch' /etc/hosts || echo '${NET_IP}.107   elasticsearch' >> /etc/hosts"
	sudo sh -c "grep -qxF '${NET_IP}.108    kibana' /etc/hosts || echo '${NET_IP}.108   kibana' >> /etc/hosts"
	sudo sh -c "grep -qxF '${NET_IP}.101    ${PROJECT_NAME}.local' /etc/hosts || echo '${NET_IP}.101   ${PROJECT_NAME}.local' >> /etc/hosts"
	sudo sh -c "grep -qxF '${NET_IP}.1   ${PROJECT_NAME}-kafka.local ui.${PROJECT_NAME}-kafka.local rest.${PROJECT_NAME}-kafka.local schema.${PROJECT_NAME}-kafka.local' /etc/hosts || echo '${NET_IP}.1   ${PROJECT_NAME}-kafka.local ui.${PROJECT_NAME}-kafka.local rest.${PROJECT_NAME}-kafka.local schema.${PROJECT_NAME}-kafka.local' >> /etc/hosts"
	make up

build:
	@echo "Building..."
	
	
	@go build -o main cmd/app/main.go

up:
	@if docker compose -f deployments/docker-compose.yml -p ${PROJECT_NAME} up -d --no-deps --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose -f deployments/docker-compose.yml -p ${PROJECT_NAME} up -d --no-deps --build; \
	fi

down:
	@if docker compose -f deployments/docker-compose.yml -p ${PROJECT_NAME} down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose -f deployments/docker-compose.yml -p ${PROJECT_NAME} down; \
	fi

clean:
	@echo "Cleaning..."
	@rm -f main

.PHONY: all build clean up down
