package types

import (
	"time"
	//	"github.com/markbates/goth"
)

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

type Session struct {
	UserID  string    `bson:"userID"`
	Token   string    `bson:"token"`
	Expires time.Time `bson:"expires"`
}

type Chat struct {
	Message string    `bson:"message" json:"message"`
	UserID  string    `bson:"userID"  json:"userID"`
	Username string
	StaffRole string
	Time    time.Time `bson:"time"    json:"time"`
}
