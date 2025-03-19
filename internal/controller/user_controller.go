package controller

import (
	"encoding/json"
	"net/http"
	"repository/internal/model"
	"repository/internal/service"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func RegisterUserRoutes(r chi.Router, userService service.UserService) {
	r.Route("/api/users", func(r chi.Router) {
		r.Get("/", listUsers(userService))
		r.Post("/", createUser(userService))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", getUserByID(userService))
			r.Put("/", updateUser(userService))
			r.Delete("/", deleteUser(userService))
		})
	})
}

// @listUsers godoc
// @Summary Список пользователей
// @Description Выводит список пользователей (можно указать ?limit=?offset=)
// @Tags Users
// @Accept json
// @Produce json
// @Param limit query int false "Сколько записей вернуть"
// @Param offset query int false "Смещение (сколько пропустить)"
// @Success 200 {array} model.User
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal error"
// @Router /api/users [get]
func listUsers(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		lim, _ := strconv.Atoi(limitStr)
		off, _ := strconv.Atoi(offsetStr)

		users, err := userService.ListUsers(ctx, lim, off)
		if err != nil {
			http.Error(w, "error get list users", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

// @createUser godoc
// @Summary Создание пользователя
// @Description Добавляет нового пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Param request body model.User true "Параметры пользователя"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {string} string "Ошибка валидации"
// @Router /api/users [post]
func createUser(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var u model.User

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := userService.CreateUser(ctx, u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": id,
		})
	}
}

// @getUserByID godoc
// @Summary Получение пользователя по ID
// @Description Выводит пользователя по переданному ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} model.User
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal error"
// @Router /api/users/{id} [get]
func getUserByID(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}

		user, err := userService.GetUserByID(ctx, id)
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// @updateUser godoc
// @Summary Обновление данных о пользователе
// @Description Обновляет данные о пользотеле на те, что передаются в запросе
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param request body model.User true "Параметры пользователя"
// @Success 200 {object} map[string]string "Успешное обновление"
// @Failure 400 {string} string "Ошибка валидации"
// @Router /api/users/{id} [put]
func updateUser(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}

		var u model.User

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u.ID = id

		err = userService.UpdateUser(ctx, u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "user upddated successfully",
		})
	}
}

// @deleteUser godoc
// @Summary Удаление пользователя по ID
// @Description Ставит метку на удаление пользователя по переданному ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]string "Пользователь успешно удален"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal error"
// @Router /api/users/{id} [delete]
func deleteUser(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}

		err = userService.DeleteUser(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "user deleted successfully",
		})
	}
}
