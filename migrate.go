package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(" Утилита управления миграциями")
		fmt.Println("================================")
		fmt.Println("Использование:")
		fmt.Println("  go run migrate.go up     - применить миграцию (создать таблицу)")
		fmt.Println("  go run migrate.go down   - откатить миграцию (удалить таблицу)")
		fmt.Println("  go run migrate.go status - показать статус миграций")
		fmt.Println("  go run migrate.go help   - показать эту помощь")
		fmt.Println("")
		fmt.Println(" Файлы миграций в папке migrations/:")
		fmt.Println("  20251212091807_create_tasks_table.up.sql")
		fmt.Println("  20251212091807_create_tasks_table.down.sql")
		return
	}

	command := os.Args[1]

	db, err := sql.Open("sqlite3", "todo.db")
	if err != nil {
		log.Fatal("Ошибка подключения к базе:", err)
	}
	defer db.Close()

	switch command {
	case "up", "migrate":
		fmt.Println(" Применение UP миграции...")
		fmt.Println(" Выполняется SQL из файла:")
		fmt.Println("   migrations/20251212091807_create_tasks_table.up.sql")
		fmt.Println("")
		fmt.Println("SQL команда:")
		fmt.Println("----------------------------------------")

		upSQL := `CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			text TEXT NOT NULL,
			is_done BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP DEFAULT NULL
		);
		
		CREATE INDEX IF NOT EXISTS idx_tasks_deleted_at ON tasks(deleted_at);`

		fmt.Println(upSQL)
		fmt.Println("----------------------------------------")

		_, err := db.Exec(upSQL)
		if err != nil {
			log.Fatal(" Ошибка выполнения SQL:", err)
		}

		fmt.Println(" Миграция применена успешно!")
		fmt.Println(" Файл базы: todo.db")
		fmt.Println(" Таблица: tasks создана")

	case "down", "rollback":
		fmt.Println(" Применение DOWN миграции...")
		fmt.Println(" Выполняется SQL из файла:")
		fmt.Println("   migrations/20251212091807_create_tasks_table.down.sql")
		fmt.Println("")
		fmt.Println("SQL команда:")
		fmt.Println("----------------------------------------")
		fmt.Println("DROP TABLE IF EXISTS tasks;")
		fmt.Println("----------------------------------------")

		_, err := db.Exec("DROP TABLE IF EXISTS tasks;")
		if err != nil {
			log.Fatal(" Ошибка выполнения SQL:", err)
		}

		fmt.Println(" Миграция откачена успешно!")
		fmt.Println(" Таблица: tasks удалена")

	case "status", "check":
		fmt.Println(" Статус миграций:")
		fmt.Println("===================")

		var tableName string
		err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='tasks'").Scan(&tableName)

		if err == sql.ErrNoRows {
			fmt.Println(" Таблица 'tasks' не существует")
			fmt.Println(" Для создания выполни: go run migrate.go up")
		} else if err != nil {
			log.Fatal("Ошибка проверки:", err)
		} else {
			fmt.Println(" Таблица 'tasks' существует")

			fmt.Println("\nСтруктура таблицы tasks:")
			fmt.Println("+------------+----------------------+-------------+")
			fmt.Println("| Поле       | Тип                  | Обязательное|")
			fmt.Println("+------------+----------------------+-------------+")

			rows, err := db.Query("PRAGMA table_info(tasks)")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			for rows.Next() {
				var cid int
				var name, typ, notnull, dflt_value string
				var pk int
				err := rows.Scan(&cid, &name, &typ, &notnull, &dflt_value, &pk)
				if err != nil {
					log.Fatal(err)
				}
				required := "НЕТ"
				if notnull == "1" {
					required = "ДА"
				}
				fmt.Printf("| %-10s | %-20s | %-11s |\n", name, typ, required)
			}
			fmt.Println("+------------+----------------------+-------------+")

			var count int
			db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
			fmt.Printf("Количество записей в таблице: %d\n", count)
		}

	case "help", "--help", "-h":
		fmt.Println(" Утилита управления миграциями")
		fmt.Println("================================")
		fmt.Println("Использование:")
		fmt.Println("  go run migrate.go up     - применить миграцию")
		fmt.Println("  go run migrate.go down   - откатить миграцию")
		fmt.Println("  go run migrate.go status - показать статус")
		fmt.Println("  go run migrate.go help   - показать эту помощь")
		fmt.Println("")
		fmt.Println(" Файлы миграций в папке migrations/:")
		fmt.Println("  20251212091807_create_tasks_table.up.sql")
		fmt.Println("  20251212091807_create_tasks_table.down.sql")

	default:
		fmt.Println(" Неизвестная команда:", command)
		fmt.Println(" Используйте: up, down, status или help")
	}
}
