package main

import (
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
)

type writerStub struct {
	err    error
	orders []Order
}

func (w *writerStub) Save(order Order) error {
	if w.err != nil {
		return w.err
	}
	w.orders = append(w.orders, order)
	return nil
}

type notifierStub struct {
	err       error
	customers []string
}

func (n *notifierStub) Send(customer string) error {
	n.customers = append(n.customers, customer)
	return n.err
}

func TestCreateOrder(t *testing.T) {
	failure := errors.New("failure")
	for _, tc := range []struct {
		name                          string
		writeErr, notifyErr           error
		wantOrders, wantNotifications int
	}{
		{"success", nil, nil, 1, 1},
		{"save failure skips notification", failure, nil, 0, 0},
		{"notification failure keeps order", nil, failure, 1, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &writerStub{err: tc.writeErr}
			notifier := &notifierStub{err: tc.notifyErr}
			err := NewOrderService(repo, notifier).CreateOrder("Иван", []string{"apple", "banana"}, 10.5)
			wantErr := tc.writeErr
			if wantErr == nil {
				wantErr = tc.notifyErr
			}
			if !errors.Is(err, wantErr) {
				t.Fatalf("error = %v, want %v", err, wantErr)
			}
			if len(repo.orders) != tc.wantOrders || len(notifier.customers) != tc.wantNotifications {
				t.Fatalf("orders = %d, notifications = %d", len(repo.orders), len(notifier.customers))
			}
			if len(repo.orders) > 0 {
				want := Order{Customer: "Иван", Products: "[apple banana]", Total: 10.5, Status: "pending"}
				if repo.orders[0] != want {
					t.Fatalf("order = %+v, want %+v", repo.orders[0], want)
				}
			}
			if len(notifier.customers) > 0 && notifier.customers[0] != "Иван" {
				t.Fatal("wrong recipient")
			}
		})
	}
}

func TestSQLiteRepository(t *testing.T) {
	db, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "orders.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := NewSQLiteRepository(db)
	for i := 0; i < 2; i++ {
		if err := repo.Init(); err != nil {
			t.Fatal(err)
		}
	}
	for _, notifier := range []Notifier{EmailSender{}, SMSSender{}} {
		if err := NewOrderService(repo, notifier).CreateOrder("O'Neil", []string{"apple"}, 10.5); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := db.Query("SELECT id, customer, products, total, status FROM orders ORDER BY id")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var got Order
		if err := rows.Scan(&got.ID, &got.Customer, &got.Products, &got.Total, &got.Status); err != nil {
			t.Fatal(err)
		}
		count++
		want := Order{ID: count, Customer: "O'Neil", Products: "[apple]", Total: 10.5, Status: "pending"}
		if got != want {
			t.Fatalf("order = %+v, want %+v", got, want)
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("count = %d, want 2", count)
	}
}

func TestMemoryRepositoryWithBothNotifiers(t *testing.T) {
	repo := &MemoryRepository{}
	for _, notifier := range []Notifier{EmailSender{}, SMSSender{}} {
		if err := NewOrderService(repo, notifier).CreateOrder("Анна", []string{"pear"}, 3); err != nil {
			t.Fatal(err)
		}
	}
	if len(repo.orders) != 2 {
		t.Fatalf("count = %d, want 2", len(repo.orders))
	}
	for i, order := range repo.orders {
		if order.ID != i+1 || order.Customer != "Анна" || order.Products != "[pear]" || order.Total != 3 || order.Status != "pending" {
			t.Fatalf("unexpected order: %+v", order)
		}
	}
}
