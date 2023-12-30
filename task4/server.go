package task4

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	ErrLessZeroNum = errors.New("id is less than zero")
	ErrTimeOut     = errors.New("timeout exceeded")
)

type Server struct {
	router chi.Router

	mu       sync.Mutex
	taskList map[int]*Task
	users    map[int]*User
}

func New(router chi.Router) *Server {
	return &Server{
		router:   router,
		taskList: make(map[int]*Task),
		users:    make(map[int]*User),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Start(addr, port string) error {
	if port == "" {
		port = "8080"
	}
	startAddr := fmt.Sprintf("%s:%s", addr, port)

	s.configureMiddleware()
	s.configureRouters()

	fmt.Println("Started")
	return http.ListenAndServe(startAddr, s)
}

func (s *Server) configureRouters() {
	s.router.Get("/tasks", s.tasks())
	s.router.Get("/tasks/{id}", s.taskInID())
	s.router.Post("/tasks/add", s.addTask())
	s.router.Delete("/tasks/{id}", s.DeleteTask())

	s.router.Put("/users", s.addUser())
	s.router.Get("/users/{id}", s.UsersByID())
	s.router.Post("/users", s.UpdateUser())
}

func (s *Server) configureMiddleware() {
	s.router.Use(s.timeoutMiddleware)
	s.router.Use(middleware.Logger)
}
func (s *Server) respond(w http.ResponseWriter, code int, data any) {
	data = map[string]any{
		"response": data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func (s *Server) error(w http.ResponseWriter, code int, err error) {
	s.respond(w, code, map[string]any{"error": err.Error()})
}
