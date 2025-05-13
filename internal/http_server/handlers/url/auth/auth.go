package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"jwt-auth-service/internal/lib/api/response"
	"jwt-auth-service/internal/lib/jwt"
	"jwt-auth-service/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	response.Response
	JWTToken string `json:"token,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=UserAuth
type UserAuth interface {
	UserExists(login, password string) (bool, error)
}

func New(log *slog.Logger, auth UserAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.auth.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		if r.Body == nil || r.ContentLength == 0 {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("empty body"))
		}

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid request body"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		if err = validator.New().Struct(req); err != nil {
			log.Error("invalid request")
			validateError := err.(validator.ValidationErrors)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateError))
			return
		}
		login, password := req.Login, req.Password
		if login == "" || password == "" {
			log.Info("empty login or password")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("empty login or password"))
			return
		}
		token, err := jwt.GenerateToken(login, password, auth)
		if err != nil {
			text := "an unexpected error occurred during token generating"
			if errors.Is(err, storage.ErrUserNotFound) {
				text = "user not found"
			}
			log.Info(text)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error(text))
			return
		}

		render.JSON(w, r, responseOk(token))
	}
}

func responseOk(token string) Response {
	return Response{
		response.OK(), token,
	}
}
