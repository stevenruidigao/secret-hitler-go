package types

import (
	"time"
)

type Game struct {
	ID                string
	Name              string
	Flag              string
	Date              time.Time
	PlayerChats       string
	PlayerCount       int
	WinningPlayerIDs  []string
	LosingPlayerIDs   []string
	WinningTeam       string
	Season            int
	Rainbow           bool
	EloMinimum        int
	Rebalance         string
	TournyFirstRound  bool
	TournySecondRound bool
	Casual            bool
	Practice          bool
	Custom            bool
	Unlisted          bool
	VerifiedOnly      bool
	Chats             []Chat
	Guesses           map[string]string
	Timer             int
	Blind             bool
	Completed         bool
}
