package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"jwt-auth-service/internal/http_server/handlers/url/auth/mocks"
	"jwt-auth-service/internal/lib/logger/handlers/slogdiscard"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthHandler(t *testing.T) {
	cases := []struct {
		name           string
		login          string
		password       string
		respError      string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Success",
			login:          "user",
			password:       "user",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Empty login",
			login:          "",
			password:       "user",
			respError:      "field Login is a required field",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty password",
			login:          "user",
			password:       "",
			respError:      "field Password is a required field",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty login and password",
			login:          "",
			password:       "",
			respError:      "field Login is a required field, field Password is a required field",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "UserExists error",
			login:          "correct_user",
			password:       "correct_password",
			respError:      "an unexpected error occurred during token generating",
			mockError:      errors.New("failed to check if user exists"),
			expectedStatus: http.StatusInternalServerError,
		},
	}
	for _, testCase := range cases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			userAuth := mocks.NewUserAuth(t)
			if tc.respError == "" || tc.mockError != nil {
				userAuth.On("UserExists", tc.login, tc.password).
					Return(false, tc.mockError).
					Once()
			}
			handler := New(slogdiscard.NewDiscardLogger(), userAuth)

			input := fmt.Sprintf(`{"login":"%s", "password":"%s"}`, tc.login, tc.password)

			req, err := http.NewRequest(http.MethodGet, "/auth", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			require.Equal(t, rr.Code, tc.expectedStatus)

			body := rr.Body.String()

			var resp Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
