.PHONY: dev backend frontend docker-up

dev:
	@echo "Starting backend..."
	cd backend && make dev &
	@echo "Starting frontend..."
	cd frontend && npm run dev

backend:
	cd backend && make dev

frontend:
	cd frontend && npm run dev

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
