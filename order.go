package main

import "fmt"

type Order struct {
	ID       int
	Customer string
	Products string
	Total    float64
	Status   string
}

// RepositoryWriter сохраняет заказ или возвращает ошибку.
type RepositoryWriter interface {
	Save(order Order) error
}

// RepositoryInitializer нужен только при подготовке хранилища.
type RepositoryInitializer interface {
	Init() error
}

// Notifier сообщает клиенту о заказе или возвращает ошибку.
// Учебные реализации имитируют отправку выводом в консоль.
type Notifier interface {
	Send(customer string) error
}

type OrderService struct {
	repository RepositoryWriter
	notifier   Notifier
}

func NewOrderService(repository RepositoryWriter, notifier Notifier) *OrderService {
	return &OrderService{repository: repository, notifier: notifier}
}

func (s *OrderService) CreateOrder(customer string, products []string, total float64) error {
	order := Order{
		Customer: customer,
		Products: fmt.Sprintf("%v", products),
		Total:    total,
		Status:   "pending",
	}
	if err := s.repository.Save(order); err != nil {
		return fmt.Errorf("сохранить заказ: %w", err)
	}
	if err := s.notifier.Send(customer); err != nil {
		return fmt.Errorf("заказ сохранён, но уведомление не отправлено: %w", err)
	}
	return nil
}
