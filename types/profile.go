package types

import (
	"time"
)

type StatSubCategory struct {
	Events    int `bson:"events"    json:"events"`
	Successes int `bson:"successes" json:"successes"`
}

type StatCategory struct {
	Liberal StatSubCategory `bson:"liberal" json:"liberal"`
	Fascist StatSubCategory `bson:"fascist" json:"fascist"`
}

type Badge struct {
	ID    string    `bson:"id" json:"id"`
	Time  time.Time `bson:"dateAwarded" json:"dateAwarded"`
	Text  string    `bson:"text" json:"text"`
	Title string    `bson:"title" json:"title"`
}

// {"stats":{"matches":{"legacyMatches":{"liberal":null,"fascist":null},"greyMatches":{"5":{"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"6":{"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"7":{"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"8":{"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"9":{"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"10":{"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"rainbowMatches":{"5":{"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"6":{"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"7":{"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"8":{"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"9":{"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"10":{"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"practiceMatches":{"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"silentMatches":{"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"emoteMatches":{"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"casualMatches":{"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"customMatches":{"liberal":{"events":0,"successes":0},"fascist":{"events":0,"successes":0}},"allMatches":{"events":309,"successes":146},"fascist":{"events":120,"successes":52},"liberal":{"events":189,"successes":94}},"actions":{"voteAccuracy":{"events":0,"successes":0},"shotAccuracy":{"events":0,"successes":0},"legacyVoteAccuracy":{"events":0,"successes":0},"legacyShotAccuracy":{"events":0,"successes":0}}},"_id":"stevengao","__v":0,"recentGames":[{"_id":"UninterestedWantingBadCat","loyalty":"liberal","playerSize":7,"isWinner":false,"date":"2021-02-11T23:49:41.470Z"},{"_id":"SilentHaplessSuccinctJackal","loyalty":"liberal","playerSize":7,"isWinner":false,"date":"2020-07-08T14:27:16.654Z"},{"_id":"OneUnknownClutteredHerringRemake","loyalty":"fascist","playerSize":7,"isWinner":true,"date":"2020-04-07T18:45:54.882Z"},{"_id":"OneUnknownClutteredHerring","loyalty":"liberal","playerSize":7,"isWinner":true,"date":"2020-04-07T18:34:43.241Z"},{"_id":"NeedlessMilitaryMaterialisticOyster","loyalty":"liberal","playerSize":7,"isWinner":false,"date":"2020-04-02T19:24:01.403Z"},{"_id":"NeighborlyDearAlcoholicMantis","loyalty":"liberal","playerSize":7,"isWinner":true,"date":"2020-04-02T19:01:10.078Z"},{"_id":"EfficaciousBouncyPaltryLobsterRemake","loyalty":"fascist","playerSize":7,"isWinner":true,"date":"2020-04-01T20:11:47.906Z"},{"_id":"EfficaciousBouncyPaltryLobster","loyalty":"fascist","playerSize":7,"isWinner":false,"date":"2020-04-01T19:58:35.617Z"},{"_id":"OldDamagedTorpidLyrebird","loyalty":"liberal","playerSize":7,"isWinner":false,"date":"2020-04-01T19:42:16.999Z"},{"_id":"ClutteredRambunctiousSameToad","loyalty":"liberal","playerSize":7,"isWinner":false,"date":"2020-04-01T16:01:56.062Z"}],"created":"03/21/2019","customCardback":"png","lastConnected":"02/22/2022","badges":[{"_id":"62149ef967f8d07cbc1dafd3","id":"eloReset1561","text":"At the time of the Elo reset, you had 1561 overall Elo and 311 games played.","title":"Elo Reset","dateAwarded":"2022-02-22T08:29:45.779Z"},{"_id":"6214e5b60a90170e98035756","id":"birthday1","text":"Your account is now 1 year old!","title":"Happy 1st birthday!","dateAwarded":"2022-02-22T13:31:34.684Z"},{"_id":"6214e5b60a90170e98035757","id":"birthday2","text":"Your account is now 2 years old!","title":"Happy 2nd birthday!","dateAwarded":"2022-02-22T13:31:34.684Z"},{"_id":"6214e5b60a90170e98035758","id":"contributor","text":"Thank you for your contributions!","title":"You contributed to the site!","dateAwarded":"2022-02-22T13:31:34.685Z"}],"eloPercentile":{},"maxElo":1600,"pastElo":[{"date":"2022-02-22T14:44:36.384Z","value":1600}],"xpOverall":311,"eloOverall":1600,"xpSeason":0,"eloSeason":1600,"isRainbowOverall":true,"isRainbowSeason":false}

type RecentGame struct {
	ID         string    `bson:"id" json:"_id"`
	Loyalty    string    `bson:"loyalty" json:"loyalty"`
	PlayerSize int       `bson:"playerSize" json:"playerSize"`
	IsWinner   bool      `bson:"isWinner" json:"isWinner"`
	Date       time.Time `bson:"date" json:"date"`
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
			StatCategory
			AllMatches StatSubCategory `bson:"allMatches" json:"allMatches"`
			// Liberal    StatSubCategory `bson:"liberal"    json:"liberal"`
			// Fascist    StatSubCategory `bson:"fascist"    json:"fascist"`
			CasualMatches StatCategory `bson:"casualMatches" json:"casualMatches"`
			CustomMatches StatCategory `bson:"customMatches" json:"customMatches"`
			EmoteMatches  StatCategory `bson:"emoteMatches" json:"emoteMatches"`
			GreyMatches   struct {
				StatCategory
				FivePlayer  StatCategory `bson:"fivePlayer" json:"5"`
				SixPlayer   StatCategory `bson:"sixPlayer" json:"6"`
				SevenPlayer StatCategory `bson:"sevenPlayer" json:"7"`
				EightPlayer StatCategory `bson:"eightPlayer" json:"8"`
				NinePlayer  StatCategory `bson:"ninePlayer" json:"9"`
				TenPlayer   StatCategory `bson:"tenPlayer" json:"10"`
			} `bson:"greyMatches" json:"greyMatches"`
			LegacyMatches   StatCategory `bson:"legacyMatches" json:"legacyMatches"`
			PracticeMatches StatCategory `bson:"practiceMatches" json:"practiceMatches"`
			RainbowMatches  struct {
				StatCategory
				FivePlayer  StatCategory `bson:"fivePlayer" json:"5"`
				SixPlayer   StatCategory `bson:"sixPlayer" json:"6"`
				SevenPlayer StatCategory `bson:"sevenPlayer" json:"7"`
				EightPlayer StatCategory `bson:"eightPlayer" json:"8"`
				NinePlayer  StatCategory `bson:"ninePlayer" json:"9"`
				TenPlayer   StatCategory `bson:"tenPlayer" json:"10"`
			} `bson:"rainbowMatches" json:"rainbowMatches"`
			SilentMatches StatCategory `bson:"silentMatches" json:"silentMatches"`
		} `bson:"matches" json:"matches"`
		Actions struct {
			VoteAccuracy       StatSubCategory `bson:"voteAccuracy" json:"voteAccuracy"`
			ShotAccuracy       StatSubCategory `bson:"shotAccuracy" json:"shotAccuracy"`
			LegacyVoteAccuracy StatSubCategory `bson:"legacyVoteAccuracy" json:"legacyVoteAccuracy"`
			LegacyShotAccuracy StatSubCategory `bson:"legacyShotAccuracy" json:"legacyShotAccuracy"`
		} `bson:"actions" json:"actions"`
	} `bson:"stats" json:"stats"`
	RecentGames    []RecentGame `bson:"recentGames" json:"recentGames"`
	Badges         []Badge      `bson:"badges" json:"badges"`
	RainbowOverall bool         `bson:"rainbowOverall" json:"isRainbowOverall"`
	RainbowSeason  bool         `bson:"rainbowSeason" json:"isRainbowSeason"`
	LastConnected  time.Time    `bson:"lastConnected" json:"lastConnected"`
	EloOverall     int          `bson:"eloOverall" json:"eloOverall"`
	EloSeason      int          `bson:"eloSeason" json:"eloSeason"`
	XPOverall      int          `bson:"xpOverall" json:"xpOverall"`
	XPSeason       int          `bson:"xpSeason" json:"xpSeason"`
}
