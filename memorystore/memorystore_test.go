package memorystore

import (
	"testing"

	"fknsrs.biz/p/requesttoken/internal/storetest"
)

func TestCreate(t *testing.T) {
	s := Store{}
	storetest.TestCreate(&s, t)
}

func TestConsume(t *testing.T) {
	s := Store{}
	storetest.TestConsume(&s, t)
}

func TestConsumeExpired(t *testing.T) {
	s := Store{}
	storetest.TestConsumeExpired(&s, t)
}

func TestConsumeWrongSession(t *testing.T) {
	s := Store{}
	storetest.TestConsumeWrongSession(&s, t)
}

func TestConsumeMissingToken(t *testing.T) {
	s := Store{}
	storetest.TestConsumeMissingToken(&s, t)
}
