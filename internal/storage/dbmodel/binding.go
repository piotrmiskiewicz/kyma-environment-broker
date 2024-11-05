package dbmodel

import (
	"time"
)

type BindingDTO struct {
	ID         string
	InstanceID string

	CreatedAt time.Time
	ExpiresAt time.Time

	Kubeconfig        string
	ExpirationSeconds int64
	CreatedBy         string
}

type BindingStats struct {
	MaxExpirationTimeInSeconds *float64 `db:"max_expiration_time_in_seconds"`
}
