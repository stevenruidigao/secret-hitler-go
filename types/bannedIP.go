package types

import (
	"time"
)

type BannedIP struct {
	Time      time.Time
	Type      string
	IP        string
	Permanent bool
}
