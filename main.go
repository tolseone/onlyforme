package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"qwerty/internal/repositories/documentRepository"
)

func main() {
	// Подключение к PostgreSQL
	db, err := sql.Open("postgres", "postgres://user:pass@localhost/dbname?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Проверка подключения
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	// Инициализация репозитория
	repo := documentRepository.NewPostgresRepository(db)

	// Далее инициализация процессора и сервера...
}
