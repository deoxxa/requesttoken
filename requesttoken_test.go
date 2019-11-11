package requesttoken

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type dummyHandler struct{}

func (dummyHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) { rw.WriteHeader(http.StatusOK) }

type dummyConverter struct {
	token []byte
}

func (c *dummyConverter) GetToken(r *http.Request) []byte {
	return c.token
}

func (c *dummyConverter) SetToken(r *http.Request, rw http.ResponseWriter, token []byte) {
	c.token = token
}

type dummyStore struct{}

func (dummyStore) Create(sessionID []byte, expiresAt time.Time) ([]byte, error) {
	return []byte("new_token"), nil
}

func (dummyStore) Consume(sessionID []byte, token []byte, now time.Time) error {
	if string(token) == "bad_session" {
		return ErrTokenSessionMismatch
	}

	if string(token) == "expired_token" {
		return ErrTokenExpired
	}

	if string(token) == "missing_token" {
		return ErrTokenNotFound
	}

	if string(token) == "good_token" {
		return nil
	}

	panic("abort")
}

func TestProvide(t *testing.T) {
	a := assert.New(t)

	var c dummyConverter

	m := Manager{Converter: &c, Store: &dummyStore{}}

	rw := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	m.Provide(dummyHandler{}).ServeHTTP(rw, r)

	a.Equal(http.StatusOK, rw.Code)
	a.Equal("new_token", string(c.token))
}

func TestEnforce(t *testing.T) {
	for _, e := range []struct {
		name   string
		token  []byte
		status int
	}{
		{"good token", []byte("good_token"), http.StatusOK},
		{"missing token", []byte("missing_token"), http.StatusForbidden},
		{"expired token", []byte("expired_token"), http.StatusForbidden},
		{"empty token", []byte(""), http.StatusForbidden},
		{"no token", nil, http.StatusForbidden},
	} {
		t.Run(e.name, func(t *testing.T) {
			a := assert.New(t)

			c := dummyConverter{token: e.token}

			m := Manager{Converter: &c, Store: &dummyStore{}}

			rw := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			m.Enforce(dummyHandler{}).ServeHTTP(rw, r)

			a.Equal(e.status, rw.Code)
		})
	}
}

func TestEnforcePanic(t *testing.T) {
	for _, e := range []struct {
		name  string
		token []byte
		err   error
	}{
		{"bad session", []byte("bad_session"), ErrTokenSessionMismatch},
		{"missing token", []byte("missing_token"), ErrTokenNotFound},
		{"expired token", []byte("expired_token"), ErrTokenExpired},
		{"empty token", []byte(""), ErrTokenNotProvided},
		{"no token", nil, ErrTokenNotProvided},
	} {
		t.Run(e.name, func(t *testing.T) {
			a := assert.New(t)

			c := dummyConverter{token: e.token}

			m := Manager{Converter: &c, Store: &dummyStore{}, Panic: true}

			rw := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			a.PanicsWithValue(e.err, func() { m.Enforce(dummyHandler{}).ServeHTTP(rw, r) })
		})
	}
}
