package jsonstorage

type Task struct {
	ID      int32  `json:"id"`
	Message string `json:"message"`
}

func NewTask(id int32, message string) *Task {
	return &Task{
		ID:      id,
		Message: message,
	}
}

type TaskList struct {
	Tasks []Task `json:"tasks"`
}
