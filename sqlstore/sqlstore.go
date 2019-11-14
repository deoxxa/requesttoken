package sqlstore

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"time"

	"fknsrs.biz/p/requesttoken"
)

const DefaultTable = "request_tokens"
const DefaultSchema = "create table if not exists request_tokens (token text not null primary key, session_id text not null, expires_at timestamp not null)"

type Store struct {
	DB    *sql.DB
	Table string
}

func (s *Store) Create(sessionID []byte, expiresAt time.Time) ([]byte, error) {
	d := make([]byte, 16)
	if _, err := rand.Read(d); err != nil {
		return nil, err
	}

	token := hex.EncodeToString(d)

	table := s.Table
	if table == "" {
		table = DefaultTable
	}

	if _, err := s.DB.Exec(
		"insert into "+table+" (token, session_id, expires_at) values ($1, $2, $3)",
		token,
		base64.StdEncoding.EncodeToString(sessionID),
		expiresAt,
	); err != nil {
		return nil, err
	}

	return []byte(token), nil
}

func (s *Store) Consume(sessionID []byte, token []byte, now time.Time) error {
	var tokenSessionID string
	var tokenExpiresAt time.Time

	table := s.Table
	if table == "" {
		table = DefaultTable
	}

	if err := s.DB.QueryRow(
		"select session_id, expires_at from "+table+" where token = $1",
		string(token),
	).Scan(&tokenSessionID, &tokenExpiresAt); err != nil {
		if err == sql.ErrNoRows {
			return requesttoken.ErrTokenNotFound
		}

		return err
	}

	d, err := base64.StdEncoding.DecodeString(tokenSessionID)
	if err != nil {
		return err
	}

	if !bytes.Equal(d, sessionID) {
		return requesttoken.ErrTokenSessionMismatch
	}

	if !tokenExpiresAt.After(now) {
		return requesttoken.ErrTokenExpired
	}

	if _, err := s.DB.Exec("delete from "+table+" where token = $1", string(token)); err != nil {
		return err
	}

	return nil
}

func (s *Store) Cleanup(now time.Time, holdExpiredRecordsFor time.Duration) {
	table := s.Table
	if table == "" {
		table = DefaultTable
	}

	_, _ = s.DB.Exec("delete from "+table+" where expires_at < $1", now.Add(0-holdExpiredRecordsFor))
}
