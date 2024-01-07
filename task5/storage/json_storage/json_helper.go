package jsonstorage

import (
	"encoding/json"
	"io"
	"os"
)

func getJsonFile(pathfile string) (*TaskList, error) {
	file, err := os.Open(pathfile)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(data) <= 0 {
		return &TaskList{
			Tasks: make([]Task, 0),
		}, nil
	}

	var taskList TaskList
	err = json.Unmarshal(data, &taskList)
	if err != nil {
		return nil, err
	}

	return &taskList, nil
}
