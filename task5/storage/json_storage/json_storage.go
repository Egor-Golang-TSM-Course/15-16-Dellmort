package jsonstorage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"
)

var (
	ErrTaskNotFound = errors.New("no tasks found under this id")
)

type Storage struct {
	file *os.File
}

func NewStorage(pathFile string) (*Storage, error) {
	file, err := os.OpenFile(pathFile, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}

	return &Storage{
		file: file,
	}, nil
}

func (s *Storage) save(ctx context.Context, message string, out chan<- int32) error {
	time.Sleep(1 * time.Second) // искуственная задержка
	defer close(out)

	tasks, err := getJsonFile(s.file.Name())
	if err != nil {
		return err
	}
	var id int32 = 1
	if tasks != nil && len(tasks.Tasks) > 0 {
		lastID := tasks.Tasks[len(tasks.Tasks)-1].ID
		if id != 0 {
			id = lastID + 1
		}
	}
	newTask := NewTask(id, message)
	tasks.Tasks = append(tasks.Tasks, *newTask)

	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	err = os.WriteFile(s.file.Name(), data, 0755)
	if err != nil {
		return err
	}

	out <- newTask.ID
	return nil
}

func (s *Storage) Save(ctx context.Context, message string) (int32, error) {
	out := make(chan int32)
	go s.save(ctx, message, out)

	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()

		case id := <-out:
			return id, nil
		}
	}
}

func (s *Storage) Get(ctx context.Context, id int32) (*Task, error) {
	tasks, err := getJsonFile(s.file.Name())
	if err != nil {
		return nil, err
	}

	if id > int32(len(tasks.Tasks)) {
		return nil, ErrTaskNotFound
	}

	task := tasks.Tasks[id-1]
	return &task, nil
}
