package types

import (
	"time"
)

type VerifyAccount struct {
	UserID     string
	Token      string
	Expiration time.Time
}
