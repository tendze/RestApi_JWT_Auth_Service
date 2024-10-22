package auth

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"jwt-auth-service/internal/lib/api/response"
	"jwt-auth-service/internal/lib/jwt"
	"jwt-auth-service/internal/storage"
	"log/slog"
	"net/http"
)

type Request struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	response.Response
	JWTToken string `json:"token,omitempty"`
}

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
			log.Error("failed to decode request body", err.Error())
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		if err = validator.New().Struct(req); err != nil {
			log.Error("invalid request", err.Error())
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
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Info("failed to generate jwt token", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("an error occurred during token generating"))
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
