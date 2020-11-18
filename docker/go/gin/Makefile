dev_up:
	docker-compose -f docker-compose-dev.yml run --name app-dev --service-ports --rm app-dev

dev_down:
	docker-compose -f docker-compose-dev.yml down

prod_up:
	docker-compose -f docker-compose-prod.yml up -d --build

prod_down:
	docker-compose -f docker-compose-prod.yml down
