all: down run

run:
	docker compose up --build

down: 
	docker compose down

run-backend:
	docker compose -f docker-compose.backend.yml up 

down-backend:
	docker compose -f docker-compose.backend.yml down 