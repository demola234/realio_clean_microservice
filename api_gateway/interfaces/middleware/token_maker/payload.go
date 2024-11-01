package token_maker

import (
	"time"
)

type Payload struct {
	Email     string    `json:"email"`
	UserID    string    `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, userID string, duration time.Duration) (*Payload, error) {

	payload := &Payload{
		Email:     username,
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}
