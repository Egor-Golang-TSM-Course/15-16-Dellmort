package task4

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

//  * r.Get("/{id}", method)
//  * articleID := chi.URLParam(r, "articleID")

// Tasks позволяет получить все задачи
func (s *Server) tasks() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			s.error(w, http.StatusRequestTimeout, ErrTimeOut)
		default:
			s.respond(w, http.StatusOK, s.taskList)
		}
	}
}

func (s *Server) addTask() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			s.error(w, http.StatusRequestTimeout, ErrTimeOut)
		default:
			var task *Task
			if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
				s.error(w, http.StatusBadRequest, err)
				return
			}

			s.mu.Lock()
			task.ID = len(s.taskList) + 1
			s.taskList[task.ID] = task

			if r.Context().Err() != nil {
				s.error(w, http.StatusRequestTimeout, ErrTimeOut)
				return
			}
			s.respond(w, http.StatusOK, fmt.Sprintf("new task id=%d", task.ID))
			s.mu.Unlock()
		}
	}
}

func (s *Server) DeleteTask() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			s.error(w, http.StatusRequestTimeout, ErrTimeOut)
		default:
			id := chi.URLParam(r, "id")
			n, err := strconv.Atoi(id)
			if err != nil {
				s.error(w, http.StatusBadRequest, errors.New("id not number"))
				return
			}
			if n < 0 {
				s.error(w, http.StatusBadRequest, ErrLessZeroNum)
				return
			}

			if r.Context().Err() != nil {
				s.error(w, http.StatusRequestTimeout, ErrTimeOut)
				return
			}

			s.mu.Lock()
			if _, ok := s.taskList[n]; !ok {
				s.error(w, http.StatusBadRequest, errors.New("task not found"))
				return
			}
			delete(s.taskList, n)
			s.mu.Unlock()

			s.respond(w, http.StatusOK, fmt.Sprintf("task %d deleted", n))
		}
	}
}

func (s *Server) taskInID() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			s.error(w, http.StatusRequestTimeout, ErrTimeOut)
		default:
			id := chi.URLParam(r, "id")
			n, err := strconv.Atoi(id)
			if err != nil {
				s.error(w, http.StatusBadRequest, errors.New("not number"))
				return
			}

			if n < 0 {
				s.error(w, http.StatusBadRequest, ErrLessZeroNum)
				return
			}

			if r.Context().Err() != nil {
				s.error(w, http.StatusRequestTimeout, ErrTimeOut)
				return
			}

			s.mu.Lock()
			s.respond(w, http.StatusOK, s.taskList[n])
			s.mu.Unlock()
		}
	}
}
