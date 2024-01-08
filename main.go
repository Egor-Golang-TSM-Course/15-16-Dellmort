package main

import (
	"lesson15_16/task5/client"
	"lesson15_16/task5/service"
	jsonstorage "lesson15_16/task5/storage/json_storage"
	"log"
	// "lesson15_16/task4"
	// "log"
	// "github.com/go-chi/chi/v5"
	// client "lesson15_16/task1"
	// "lesson15_16/task2"
	// "lesson15_16/task3"
)

func main() {
	// client.Start()
	// task2.Start()
	// task3.Start()

	// chi := chi.NewRouter()
	// server := task4.New(chi)
	// if err := server.Start("localhost", "8080"); err != nil {
	// 	log.Fatal(err)
	// }]

	go client.StartClient()

	storage, err := jsonstorage.NewStorage("task5/tasks/tasks.json")
	if err != nil {
		log.Fatal(err)
	}

	s, _ := service.NewGRPCStorage(storage)
	s.Start()
}
