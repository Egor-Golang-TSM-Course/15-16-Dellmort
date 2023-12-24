package task3

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	SMS   = "SMS"
	EMAIL = "EMAIL"
	PUSH  = "PUSH"
)

type Notification struct {
	ctx   context.Context
	sms   chan string
	email chan string
	push  chan string
}

func NewNotification(ctx context.Context) *Notification {
	return &Notification{
		ctx:   ctx,
		sms:   make(chan string),
		email: make(chan string),
		push:  make(chan string),
	}
}

func (n *Notification) Process() {
	ctx, cancel := context.WithCancel(n.ctx)
	defer cancel()
	n.ctx = ctx
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Завершаю работу")
			close(n.sms)
			close(n.email)
			close(n.push)
			os.Exit(0)
		case <-sigChan:
			cancel()
		case message, ok := <-n.email:
			if ok {
				fmt.Println("Принял email уведомление:", message)
			}
		case message, ok := <-n.sms:
			if ok {
				fmt.Println("Принял sms уведомление:", message)
			}
		case message, ok := <-n.push:
			if ok {
				fmt.Println("Принял push уведомление:", message)
			}
		}
	}
}

func (n *Notification) Send(notificType string, message string) {
	for {
		select {
		case <-n.ctx.Done():
			fmt.Println("Выхожу по контексту")
			return
		default:
			switch notificType {
			case SMS:
				n.sms <- message
			case EMAIL:
				n.email <- message
			case PUSH:
				n.push <- message
			default:
				fmt.Println("undefined type")
			}
			return
		}
	}
}

func Start() {
	ctx := context.Background()
	n := NewNotification(ctx)
	go n.Process()
	go n.Send(SMS, "HELLO")
	go n.Send(PUSH, "message")
	go n.Send(EMAIL, "happy")
	go n.Send(SMS, "HELLO")
	go n.Send(PUSH, "HELLO")
	go n.Send(SMS, "HELLO")

	time.Sleep(5 * time.Second)
}
