# Поднимаем контейнер PostgreSQL в docker-compose
up:
	docker-compose up -d postgres

# Останавливаем контейнер PostgreSQL в docker-compose
down:
	docker-compose down

.PHONY: \
  up \
  down