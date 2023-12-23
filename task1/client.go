package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	return &Client{
		client: http.DefaultClient,
	}
}

func (c *Client) Request(ctx context.Context, method string, url string, timeout time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	requestCh := make(chan *http.Request)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	go func() {
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			log.Fatal(err)
		}
		requestCh <- req
		close(requestCh)
	}()

	for loop := true; loop; {
		select {
		case <-ctx.Done():
			fmt.Println("Выхожу по таймауту")
			return
		case req, ok := <-requestCh:
			if !ok {
				continue
			}

			resp, err := c.client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(body))
			loop = false
		}
	}
}

func Start() {
	var wg sync.WaitGroup
	ctx := context.Background()
	client := NewClient()

	wg.Add(3)
	go client.Request(ctx, http.MethodGet, "http://google.com", 1*time.Second, &wg)
	go client.Request(ctx, http.MethodGet, "https://youtube.com", 1*time.Second, &wg)
	go client.Request(ctx, http.MethodGet, "https://github.com", 1*time.Second, &wg)

	wg.Wait()
}
