.PHONY: help install run-backend run-frontend run docker-up docker-down clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install all dependencies
	@echo "Installing backend dependencies..."
	cd backend && go mod download
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

run-backend: ## Run backend server
	@echo "Starting backend server..."
	cd backend && go run main.go

run-frontend: ## Run frontend development server
	@echo "Starting frontend server..."
	cd frontend && npm start

run: ## Run both backend and frontend (requires two terminals)
	@echo "Please run 'make run-backend' in one terminal and 'make run-frontend' in another"

docker-up: ## Start all services with Docker Compose
	docker-compose up -d
	@echo "Services started! Backend: http://localhost:8080, Frontend: http://localhost:3000"

docker-down: ## Stop all Docker services
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

clean: ## Clean build artifacts
	cd backend && go clean
	cd frontend && rm -rf node_modules build
	rm -rf backend/uploads

setup-env: ## Create .env files from examples
	cp backend/.env.example backend/.env
	cp frontend/.env.example frontend/.env
	@echo "Environment files created. Please update them with your settings."
