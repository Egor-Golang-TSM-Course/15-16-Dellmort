package task2

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func generator(ctx context.Context, data chan int) {
	go process(ctx, data)
	for {
		select {
		case <-ctx.Done():
			return
		case data <- rand.Intn(555):
			fmt.Println("Сгенерировал значение")
			return
		}
	}
}

func process(ctx context.Context, data chan int) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		case procData, ok := <-data:
			if !ok {
				return
			}
			go finally(ctx, data)

			// какая-то обработка
			data <- procData * procData * 2
			close(data)
			return
		}
	}
}

func finally(ctx context.Context, data <-chan int) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case value, ok := <-data:
			if !ok {
				return
			}
			fmt.Println("finally result: ", value)
			return
		}
	}
}

func Start() {
	data := make(chan int)
	ctx := context.Background()
	generator(ctx, data)
	time.Sleep(1 * time.Second)
}
