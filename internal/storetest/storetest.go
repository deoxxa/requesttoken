package storetest

import (
	"testing"
	"time"

	"fknsrs.biz/p/requesttoken"

	"github.com/stretchr/testify/assert"
)

var (
	session1 = []byte("session1")
	session2 = []byte("session2")
	session3 = []byte("session3")
)

func TestCreate(store requesttoken.Store, t *testing.T) {
	a := assert.New(t)

	token, err := store.Create(session1, time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC))
	a.NoError(err)
	a.NotNil(token)
	a.NotEmpty(token)
}

func TestConsume(store requesttoken.Store, t *testing.T) {
	a := assert.New(t)

	token, err := store.Create(session1, time.Date(2019, time.January, 1, 1, 0, 0, 0, time.UTC))
	a.NoError(err)
	a.NotNil(token)
	a.NotEmpty(token)

	a.NoError(store.Consume(session1, token, time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)))
}

func TestConsumeExpired(store requesttoken.Store, t *testing.T) {
	a := assert.New(t)

	token, err := store.Create(session1, time.Date(2019, time.January, 1, 1, 0, 0, 0, time.UTC))
	a.NoError(err)
	a.NotNil(token)
	a.NotEmpty(token)

	a.EqualError(store.Consume(session1, token, time.Date(2019, time.January, 2, 0, 0, 0, 0, time.UTC)), requesttoken.ErrTokenExpired.Error())
}

func TestConsumeWrongSession(store requesttoken.Store, t *testing.T) {
	a := assert.New(t)

	token, err := store.Create(session1, time.Date(2019, time.January, 1, 1, 0, 0, 0, time.UTC))
	a.NoError(err)
	a.NotNil(token)
	a.NotEmpty(token)

	a.EqualError(store.Consume(session2, token, time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)), requesttoken.ErrTokenSessionMismatch.Error())
}

func TestConsumeMissingToken(store requesttoken.Store, t *testing.T) {
	a := assert.New(t)

	a.EqualError(store.Consume(session1, []byte("missing"), time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)), requesttoken.ErrTokenNotFound.Error())
}
