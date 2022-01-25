package types

import (
	"time"
)

type Report struct {
	Time              time.Time
	Game              string
	ReportedPlayerID  string
	Reason            string
	ReportingPlayerID string
	GameType          string
	Comment           string
	Active            bool
}
