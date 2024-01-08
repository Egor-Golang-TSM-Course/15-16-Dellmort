package task4

type Task struct {
	ID      int
	Message string `json:"message"`
}

type User struct {
	Name string
	Age  int
}
