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

type GeneralChat struct {
	Message   string    `bson:"message" json:"chat"`
	UserID    string    `bson:"userID"  json:"userID"`
	Username  string    `bson:"username" json:"userName"`
	StaffRole string    `bson:"staffRole" json:"staffRole"`
	Timestamp time.Time `bson:"timestamp"    json:"timestamp"`
}

type GeneralChats struct {
	List   []GeneralChat `bson:"list" json:"list"`
	Sticky string        `bson:"sticky" json:"sticky"`
}

type PlayerChat struct {
	Message   string    `bson:"message" json:"chat"`
	UserID    string    `bson:"userID"  json:"userID"`
	Username  string    `bson:"username" json:"userName"`
	StaffRole string    `bson:"staffRole" json:"staffRole"`
	Timestamp time.Time `bson:"timestamp"    json:"timestamp"`
	GameID    string    `bson:"gameID" json:"uid"`
}

type GameChat struct {
	Text string `bson:"text" json:"text"`
	Type string `bson:"type" json:"type"`
}
