package apiserver

import (
	"fmt"
	"goodisgood/internal/app/model"
	"goodisgood/internal/app/storage/teststorage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

func TestServerAuthenticateUser(t *testing.T) {
	storage := teststorage.NewStorage()
	a := model.TestAccount()
	storage.Account().Create(a)

	testCases := []struct {
		name         string
		cookie       map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "authenticated",
			cookie: map[interface{}]interface{}{
				"account_id": a.UUID,
			},
			expectedCode: http.StatusOK,
		},
		// {
		// 	name:         "not authenticated",
		// 	cookie:       nil,
		// 	expectedCode: http.StatusUnauthorized,
		// },
	}

	secret := []byte("secret")
	s := newServer(storage, sessions.NewCookieStore(secret))
	sc := securecookie.New(secret, nil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			cookieStr, _ := sc.Encode(sessionName, tc.cookie)
			req.Header.Set("Cookie ", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			s.authenticateAccount(handler).ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
