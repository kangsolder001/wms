.PHONY: dev backend frontend docker-up docker-prod docker-prod-down

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

docker-prod:
	docker compose -f docker-compose.prod.yml up --build -d

docker-prod-down:
	docker compose -f docker-compose.prod.yml down
