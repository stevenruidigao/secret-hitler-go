package types

import (
	"time"
)

type ModAction struct {
	Time        time.Time
	ModeratorID string
	IP          string
	UserID      string
	Notes       string
	Action      string
}
