package types

import (
	"time"
)

type ProfileCategory struct {
	Events    int `bson:"events"    json:"events"`
	Successes int `bson:"successes" json:"successes"`
}

type Profile struct {
	UserID          string    `bson:"userID"          json:"userID"`
	Username        string    `bson:"username"        json:"username"`
	Version         string    `bson:"version"         json:"version"` // versioning for `recalculateProfiles`
	Created         time.Time `bson:"created"         json:"created"`
	CustomCardback  string    `bson:"cardback"        json:"cardback"`
	Bio             string    `bson:"bio"             json:"bio"`
	LastConnectedIP string    `bson:"lastConnectedIP" json:"lastConnectedIP"`
	Stats           struct {
		Matches struct {
			AllMatches ProfileCategory `bson:"allMatches" json:"allMatches"`
			Liberal    ProfileCategory `bson:"liberal"    json:"liberal"`
			Fascist    ProfileCategory `bson:"fascist"     json:"fascist"`
		} `bson:"matches" json:"matches"`
		Actions struct {
			VoteAccuracy ProfileCategory `bson:"voteAccuracy" json:"voteAccuracy"`
			ShotAccuracy ProfileCategory `bson:"shotAccuracy" json:"shotAccuracy"`
		} `bson:"actions" json:"actions"`
	} `bson:"stats" json:"stats"`
	RecentGames []GamePublic `bson:"recentGames" json:"recentGames"`
}
