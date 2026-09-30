package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	db, err := sql.Open("sqlite3", "orders.db")
	if err != nil {
		return err
	}
	defer db.Close()

	repository := NewSQLiteRepository(db)
	var initializer RepositoryInitializer = repository
	if err := initializer.Init(); err != nil {
		return err
	}

	emailService := NewOrderService(repository, EmailSender{})
	if err := emailService.CreateOrder("Иван", []string{"apple", "banana"}, 10.5); err != nil {
		return err
	}

	smsService := NewOrderService(repository, SMSSender{})
	if err := smsService.CreateOrder("Мария", []string{"orange"}, 5); err != nil {
		return err
	}

	// Замена хранилища также не требует изменения OrderService.
	memoryService := NewOrderService(&MemoryRepository{}, EmailSender{})
	return memoryService.CreateOrder("Анна", []string{"pear"}, 3)
}
