package memorystore

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"time"

	"fknsrs.biz/p/requesttoken"
)

type record struct {
	sessionID []byte
	token     []byte
	expiry    time.Time
}

type Store struct {
	records []record
}

func (s *Store) Create(sessionID []byte, expiresAt time.Time) ([]byte, error) {
	d := make([]byte, 16)
	if _, err := rand.Read(d); err != nil {
		return nil, err
	}

	token := []byte(hex.EncodeToString(d))

	s.records = append(s.records, record{sessionID, token, expiresAt})

	return token, nil
}

func (s *Store) Consume(sessionID []byte, token []byte, now time.Time) error {
	for i := 0; i < len(s.records); i++ {
		e := s.records[i]

		if !bytes.Equal(e.token, token) {
			continue
		}

		s.records[i] = s.records[len(s.records)-1]
		s.records = s.records[:len(s.records)-1]

		if !e.expiry.After(now) {
			return requesttoken.ErrTokenExpired
		}

		if !bytes.Equal(e.sessionID, sessionID) {
			return requesttoken.ErrTokenSessionMismatch
		}

		return nil
	}

	return requesttoken.ErrTokenNotFound
}

func (s *Store) Cleanup(now time.Time, holdExpiredRecordsFor time.Duration) {
	for i := 0; i < len(s.records); i++ {
		e := s.records[i]

		if now.Sub(e.expiry) < holdExpiredRecordsFor {
			continue
		}

		s.records[i] = s.records[len(s.records)-1]
		s.records = s.records[:len(s.records)-1]
	}
}

func (s *Store) Length() int { return len(s.records) }
