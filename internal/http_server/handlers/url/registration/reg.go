package registration

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"jwt-auth-service/internal/storage"
	"jwt-auth-service/lib/api/response"
	"log/slog"
	"net/http"
)

type Request struct {
	Login    string `json:"login" validate:"required,login"`
	Password string `json:"password" validate:"required,password"`
}

type Response struct {
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
	JWTToken string `json:"jwt-token,omitempty"`
}

type USERRegister interface {
	SaveUser(login, password string) (int64, error)
}

func New(log *slog.Logger, register USERRegister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.registration.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
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
			render.JSON(w, r, response.ValidationError(validateError))
			return
		}
		login, password := req.Login, req.Password
		if login == "" || password == "" {
			log.Info("empty login or password")
			render.JSON(w, r, response.Error("empty login or password"))
			return
		}
		lastId, err := register.SaveUser(login, password)
		if errors.Is(err, storage.ErrUserExists) {
			log.Error("user %v already exists", err.Error())
			render.JSON(w, r, response.Error("user already exists"))
			return
		}
		if err != nil {
			log.Error("failed to save user", err.Error())
			render.JSON(w, r, response.Error("failed to svae user"))
			return
		}
		log.Info(fmt.Sprintf("user: <%v> successfully saved with id %v", login, lastId))

	}
}
