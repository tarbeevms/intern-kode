package middleware

import (
	"context"
	"myapp/internal/logic"
	"myapp/pkg/io"
	"net/http"
	"strings"
)

var (
	// Список URL, для которых не требуется аутентификация
	noAuthUrls = map[string]struct{}{
		"/api/login": {},
	}
)

// Authenticated middleware проверяет JWT и наличие сессии в базе данных.
func Authenticated(next http.Handler, ll *logic.LogicLayer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверка, есть ли текущий URL в списке исключений noAuthUrls
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			// Если URL не требует авторизации, передаем управление следующему обработчику
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")

		// Проверка на правильный формат заголовка Authorization
		if len(t) != 2 {
			io.SendError(w, "Not authorized, wrong header format", http.StatusUnauthorized)
			return
		}
		if t[0] != "Bearer" {
			io.SendError(w, "Not authorized, wrong header format", http.StatusUnauthorized)
			return
		}

		authToken := t[1]

		// Проверяем авторизацию и существование сессии в базе данных, не просрочена ли она
		authorized, err := ll.IsAuthorized(authToken)
		if err != nil {
			strError := "Not authorized: " + err.Error()
			io.SendError(w, strError, http.StatusUnauthorized)
			return
		} else if !authorized {
			io.SendError(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		// Добавляем токен в контекст
		ctx := context.WithValue(r.Context(), "token", authToken)

		// Передаем управление следующему обработчику с обновленным контекстом
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
