package main

import (
	"database/sql"
	"sync"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) Init() error {
	_, err := r.db.Exec(`CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		customer TEXT NOT NULL,
		products TEXT NOT NULL,
		total REAL NOT NULL,
		status TEXT NOT NULL
	)`)
	return err
}

func (r *SQLiteRepository) Save(order Order) error {
	_, err := r.db.Exec(
		"INSERT INTO orders (customer, products, total, status) VALUES (?, ?, ?, ?)",
		order.Customer, order.Products, order.Total, order.Status,
	)
	return err
}

// MemoryRepository готов к использованию без создания таблиц.
type MemoryRepository struct {
	mu     sync.Mutex
	orders []Order
}

func (r *MemoryRepository) Save(order Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	order.ID = len(r.orders) + 1
	r.orders = append(r.orders, order)
	return nil
}

var _ RepositoryWriter = (*SQLiteRepository)(nil)
var _ RepositoryInitializer = (*SQLiteRepository)(nil)
var _ RepositoryWriter = (*MemoryRepository)(nil)
