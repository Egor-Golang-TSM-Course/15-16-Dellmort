package task4

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"
)

func (s *Server) addUser() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			fmt.Println("exit")
			s.error(w, http.StatusRequestTimeout, ErrTimeOut)
			return
		default:
			name := r.FormValue("name")
			age := r.FormValue("age")
			if name == "" {
				s.error(w, http.StatusBadRequest, errors.New("name is null"))
				return
			}
			n, err := strconv.Atoi(age)
			if err != nil {
				s.error(w, http.StatusBadRequest, err)
				return
			}

			user := &User{
				Name: name,
				Age:  n,
			}
			if r.Context().Err() != nil {
				s.error(w, http.StatusRequestTimeout, ErrTimeOut)
				return
			}

			s.mu.Lock()
			s.users[len(s.users)+1] = user
			s.mu.Unlock()
			s.respond(w, http.StatusOK, fmt.Sprintf("ok! new user = %v", user))
		}
	}
}

func (s *Server) UsersByID() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			s.error(w, http.StatusRequestTimeout, ErrTimeOut)
		default:
			id := chi.URLParam(r, "id")
			n, err := strconv.Atoi(id)
			if err != nil {
				s.error(w, http.StatusBadRequest, err)
			}
			if n < 0 {
				s.error(w, http.StatusBadRequest, ErrLessZeroNum)
			}
			if r.Context().Err() != nil {
				s.error(w, http.StatusRequestTimeout, ErrTimeOut)
				return
			}

			s.mu.Lock()
			user := s.users[n]
			s.mu.Unlock()

			s.respond(w, http.StatusOK, user)
		}
	}
}

// post
func (s *Server) UpdateUser() func(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID   int    `json:"id"`
		Name string `json:"name,omitempty"`
		Age  int    `json:"age,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			s.error(w, http.StatusRequestTimeout, ErrTimeOut)
			return
		default:
			var resp response
			err := json.NewDecoder(r.Body).Decode(&resp)
			if err != nil {
				s.error(w, http.StatusInternalServerError, err)
				return
			}

			if resp.ID <= 0 {
				s.error(w, http.StatusBadRequest, ErrLessZeroNum)
				return
			}

			s.mu.Lock()
			defer s.mu.Unlock()
			if user, ok := s.users[resp.ID]; ok {
				if user.Name == "" && utf8.RuneCountInString(resp.Name) <= 3 {
					s.error(w, http.StatusBadRequest, errors.New("name must be more than 3 characters"))
					return
				}
				if resp.Age <= 0 {
					s.error(w, http.StatusBadRequest, errors.New("возраст меньше или равен нулю"))
				}
				user.Name = resp.Name
				user.Age = resp.Age

				if r.Context().Err() != nil {
					s.error(w, http.StatusRequestTimeout, ErrTimeOut)
					return
				}

				s.respond(w, http.StatusOK, s.users[resp.ID])
				return
			}

			s.error(w, http.StatusBadRequest, errors.New("данного юзера не существует"))
		}
	}
}
