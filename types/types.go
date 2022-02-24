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
	UserID  string    `bson:"userID"  json:"userID`
	Token   string    `bson:"token"   json:"token"`
	Expires time.Time `bson:"expires" json:"expires"`
}

type GeneralChat struct {
	Message   string    `bson:"message"   json:"chat"`
	UserID    string    `bson:"userID"    json:"userID"`
	Username  string    `bson:"username"  json:"userName"`
	StaffRole string    `bson:"staffRole" json:"staffRole"`
	Timestamp time.Time `bson:"timestamp" json:"time"`
}

type GeneralChats struct {
	List   []GeneralChat `bson:"list"   json:"list"`
	Sticky string        `bson:"sticky" json:"sticky"`
}
