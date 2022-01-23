package types

import (
	"time"

	"github.com/markbates/goth"
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

type GameSettings struct {
	CustomWidth string
	FontFamily  string
}

type User struct {
	goth.User
	Cardback     string    `bson:"cardback"`
	LocalUserID  string    `bson:"localUserID"`
	Sessions     []Session `bson:"sessions"`
	Username     string    `bson:"username"`
	GameSettings GameSettings
}
