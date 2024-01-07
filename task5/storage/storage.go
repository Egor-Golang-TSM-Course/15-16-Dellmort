package storage

import (
	"context"
	jsonstorage "lesson15_16/task5/storage/json_storage"
)

type Storage interface {
	Save(context.Context, string) (int32, error)
	Get(context.Context, int32) (*jsonstorage.Task, error)
}
