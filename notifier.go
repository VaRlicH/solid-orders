package main

import "fmt"

type EmailSender struct{}

func (EmailSender) Send(customer string) error {
	_, err := fmt.Printf("[Email, имитация] Уведомление клиенту %s\n", customer)
	return err
}

type SMSSender struct{}

func (SMSSender) Send(customer string) error {
	_, err := fmt.Printf("[SMS, имитация] Уведомление клиенту %s\n", customer)
	return err
}

var _ Notifier = EmailSender{}
var _ Notifier = SMSSender{}
