package validate

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"jwt-auth-service/internal/lib/api/response"
	"jwt-auth-service/internal/lib/jwt"
	"log/slog"
	"net/http"
)

type Response struct {
	response.Response
	UserLogin string `json:"user-login"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.validate.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Info("Header missing")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("Header missing"))
			return
		}

		var tokenString string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenString)
		if err != nil {
			log.Info("Invalid authorization header format")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("Invalid authorization header format"))
			return
		}
		userLogin, err := jwt.ValidateToken(tokenString)
		if err != nil {
			log.Info("Invalid token")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("Invalid token"))
			return
		}
		render.JSON(w, r, responseOK(userLogin))
	}
}

func responseOK(userLogin string) Response {
	return Response{
		response.OK(),
		userLogin,
	}
}
