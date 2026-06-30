package models

import (
	"time"

	"github.com/google/uuid"
)

type Activities struct {
	ID         uuid.UUID
	user_id    uuid.UUID
	weight     int
	logged_at  time.Time
	created_at time.Time
}
