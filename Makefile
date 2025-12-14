# Простой Makefile для SQLite миграций

DB_DSN := "sqlite3://todo.db"
MIGRATE := migrate -path ./migrations -database $(DB_DSN)

# Создать новую миграцию
migrate-new:
	migrate create -ext sql -dir ./migrations $(NAME)

# Применить миграции
migrate-up:
	@echo "🔼 Применяю миграции..."
	$(MIGRATE) up

# Откатить миграции
migrate-down:
	@echo "🔽 Откатываю миграции..."
	$(MIGRATE) down

# Показать статус
migrate-status:
	@echo "📊 Статус миграций:"
	$(MIGRATE) version

# Запуск приложения
run:
	@echo "🚀 Запускаю приложение..."
	go run cmd/app/main.go

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  make migrate-new NAME=tasks - создать миграцию"
	@echo "  make migrate-up             - применить миграции"
	@echo "  make migrate-down           - откатить миграции"
	@echo "  make migrate-status         - показать статус"
	@echo "  make run                    - запустить приложение"
	@echo "  make help                   - эта справка"