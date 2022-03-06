package types

import (
	"time"
)

type Report struct {
	Time                   time.Time `bson:"time" json:"time"`
	GameID                 string    `bson:"gameID" json:"gameID"`
	ReportedPlayerID       string    `bson:"reportedPlayerID" json:"reportedPlayerID"`
	Reason                 string    `bson:"reason" json:"reason"`
	ReportingPlayerID      string    `bson:"reportingPlayerID" json:"reportingPlayerID"`
	GameType               string    `bson:"gameType" json:"gameType"`
	Comment                string    `bson:"comment" json:"comment"`
	Active                 bool      `bson:"active" json:"active"`
	ReportedPlayerUsername string    `bson:"reportedPlayerUsername" json:"reportedPlayerUsername"`
	ReportedPlayerSeat     int       `bson:"reportedPlayerSeat" json:"reportedPlayerSeat"`
	ReportedPlayerRole     string    `bson:"reportedPlayerRole" json:"reportedPlayerRole"`
	GameName               string    `bson:"gameName" json:"gameName"`
	GameElectionCount      int       `bson:"gameElectionCount" json:"gameElectionCount"`
	ReportType             string    `bson:"reportType" json:"reportType"`
}
