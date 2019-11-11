package requesttoken

import (
	"errors"
	"net/http"
	"time"
)

var (
	ErrTokenExpired         = errors.New("token has expired")
	ErrTokenNotFound        = errors.New("token not found")
	ErrTokenSessionMismatch = errors.New("token was for another session")
	ErrTokenNotProvided     = errors.New("token not provided")
)

var (
	DummySessionID = []byte("dummy_session")
)

type Converter interface {
	GetToken(r *http.Request) []byte
	SetToken(r *http.Request, rw http.ResponseWriter, token []byte)
}

type ConverterSession interface {
	GetSessionAndToken(r *http.Request) ([]byte, []byte)
}

type Store interface {
	Create(sessionID []byte, expiresAt time.Time) ([]byte, error)
	Consume(sessionID []byte, token []byte, now time.Time) error
}

type Manager struct {
	Converter Converter
	Store     Store
	TTL       time.Duration
	Panic     bool
}

func getSessionAndToken(c Converter, r *http.Request) ([]byte, []byte) {
	if cs, ok := c.(ConverterSession); ok {
		return cs.GetSessionAndToken(r)
	}

	return DummySessionID, c.GetToken(r)
}

func (m *Manager) Provide(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if sessionID, _ := getSessionAndToken(m.Converter, r); sessionID != nil {
			token, err := m.Store.Create(sessionID, time.Now().Add(m.TTL))
			if err != nil {
				panic(err)
			}

			m.Converter.SetToken(r, rw, token)
		}

		handler.ServeHTTP(rw, r)
	})
}

func (m *Manager) Enforce(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		sessionID, tokenID := getSessionAndToken(m.Converter, r)

		var err error

		if tokenID == nil || len(tokenID) == 0 {
			err = ErrTokenNotProvided
		} else {
			err = m.Store.Consume(sessionID, tokenID, time.Now())
		}

		if err != nil {
			if m.Panic {
				panic(err)
			}

			switch err {
			case ErrTokenExpired, ErrTokenNotFound, ErrTokenSessionMismatch, ErrTokenNotProvided:
				http.Error(rw, "", http.StatusForbidden)
				return
			}

			http.Error(rw, "", http.StatusInternalServerError)
		}

		handler.ServeHTTP(rw, r)
	})
}
