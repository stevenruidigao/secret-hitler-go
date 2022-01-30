package types

import (
	"time"
)

type DeckState struct {
	Liberal int `bson:"lib" json:"lib"`
	Fascist int `bson:"fas" json:"fas"`
}

/*type TrackState struct {
	Liberal int `bson:"lib" json:"lib"`
	Fascist int `bson:"fas" json:"fas"`
}*/

type CustomGameSettings struct {
	Enabled bool     `bson:"enabled" json:"enabled"`
	Powers  []string `bson:"powers"  json:"powers"`
}

type GameState struct {
	PreviousElectedGovernment []int `bson:"previousElectedGovernment" json:"previousElectedGovernment"`
	UndrawnPolicyCount        int   `bson:"undrawnPolicyCount"        json:"undrawnPolicyCount"`
	DiscardedPolicyCount      int   `bson:"discardedPolicyCount"      json:"discardedPolicyCount"`
	PresidentIndex            int   `bson:"presidentIndex"            json:"presidentIndex"`
}

type TrackState struct {
	LiberalPolicyCount   int      `bson:"liberalPolicyCount"   json:"liberalPolicyCount"`
	FascistPolicyCount   int      `bson:"fascistPolicyCount"   json:"fascistPolicyCount"`
	ElectionTrackerCount int      `bson:"electionTrackerCount" json:"electionTrackerCount"`
	EnactedPolicies      []string `bson:"enactedPolicies"      json:"enactedPolicies"`
}

type GeneralGameSettings struct {
	WhitelistedPlayers      []string  `bson:"whitelistedPlayers"      json:"whitelistedPlayers"`
	ID                      string    `bson:"id"                      json:"uid"`
	Name                    string    `bson:"name"                    json:"name"`
	Flag                    string    `bson:"flag"                    json:"flag"`
	MinPlayersCount         int       `bson:"minPlayersCount"         json:"minPlayersCount"`
	ExcludedPlayerCount     []int     `bson:"excludedPlayerCount"     json:"excludedPlayerCount"`
	MaxPlayersCount         int       `bson:"maxPlayersCount"         json:"maxPlayersCount"`
	Status                  string    `bson:"status"                  json:"status"`
	Experienced             bool      `bson:"experienced"             json:"experiencedMode"`
	PlayerChats             string    `bson:"playerChats"             json:"playerChats"`
	VerifiedOnly            bool      `bson:"verifiedOnly"            json:"isVerifiedOnly"`
	DisableObserverLobby    bool      `bson:"disableObserverLobby"           json:"disableObserverLobby"`
	DisableObserver         bool      `bson:"disableObserver"                json:"disableObserver"`
	Tourny                  bool      `bson:"tourny"                  json:"isTourny"`
	LastModPing             int       `bson:"lastModPing"             json:"lastModPing"`
	ChatReplTime            []int     `bson:"chatReplTime"            json:"chatReplTime"`
	DisableGamechat         bool      `bson:"disableGamechat"                json:"disableGamechat"`
	Rainbow                 bool      `bson:"rainbow"                 json:"rainbowgame"`
	Blind                   bool      `bson:"blind"                   json:"blindMode"`
	Timer                   int       `bson:"timer"           	     json:"timedMode"`
	Flappy                  bool      `bson:"flappy"                  json:"flappyMode"`
	FlappyOnly              bool      `bson:"flappyOnly"              json:"flappyOnlyMode"`
	Casual                  bool      `bson:"casual"                  json:"casualGame"`
	Practice                bool      `bson:"practice"                json:"practiceGame"`
	Rebalance6p             bool      `bson:"rebalance6p"             json:"rebalance6p"`
	Rebalance7p             bool      `bson:"rebalance7p"             json:"rebalance7p"`
	Rebalance9p2f           bool      `bson:"rebalance9p2f"           json:"rebalance9p2f"`
	Unlisted                bool      `bson:"unlisted"                json:"unlistedGame"`
	Private                 bool      `bson:"private"                 json:"private"`
	PrivatePassword         string    `bson:"privatePassword"         json:"privatePassword"`
	PrivateAnonymousRemakes bool      `bson:"privateAnonymousRemakes" json:"privateAnonymousRemakes"`
	PrivateOnly             bool      `bson:"privateOnly"             json:"privateOnly"`
	ElectionCount           int       `bson:"electionCount"           json:"electionCount"`
	Remade                  bool      `bson:"remade"                  json:"isRemade"`
	EloMinimum              int       `bson:"eloMinimum"              json:"eloSliderValue"`
	EloMaximum              int       `bson:"eloMaximum"              json:"eloMaximum"`
	TimeCreated             time.Time `bson:"timeCreated" json:"timeCreated"`
	Usernames               []string  `bson:"usernames" json:"userNames"`
	CustomCardback          []string  `bson:"customCardback" json:"customCardback"`
	CustomCardbackUID       []string  `bson:"customCardbackUID" json:"customCardbackUid"`
	Players                 []User    `bson:"players" json:"players"`
	SeatedCount             int       `bson:"seatedCount" json:"seatedCount"`
}

type Game struct {
	ID                  string                 `bson:"id"                      json:"uid"`
	Name                string                 `bson:"name"                    json:"name"`
	Flag                string                 `bson:"flag"                    json:"flag"`
	Date                time.Time              `bson:"date"                    json:"date"`
	PlayerChats         string                 `bson:"playerChats"             json:"playerChats"`
	PlayerCount         int                    `bson:"playerCount"             json:"playerCount"`
	WinningPlayerIDs    []string               `bson:"winningPlayers"          json:"winningPlayers"`
	LosingPlayerIDs     []string               `bson:"losingPlayers"           json:"losingPlayers"`
	WinningTeam         string                 `bson:"winningTeam"             json:"winningTeam"`
	Season              int                    `bson:"season"                  json:"season"`
	Rainbow             bool                   `bson:"rainbow"                 json:"rainbowgame"`
	EloMinimum          int                    `bson:"eloMinimum"              json:"eloMinimum"`
	EloMaximum          int                    `bson:"eloMaximum"              json:"eloMaximum"`
	Rebalanced          string                 `bson:"rebalanced"              json:"rebalanced"`
	Rebalance6p         bool                   `bson:"rebalance6p"             json:"rebalance6p"`
	Rebalance7p         bool                   `bson:"rebalance7p"             json:"rebalance7p"`
	Rebalance9p         bool                   `bson:"rebalance9p"             json:"rebalance9p"`
	Rebalance9p2f       bool                   `bson:"rebalance9p2f"           json:"rebalance9p2f"`
	TournyFirstRound    bool                   `bson:"tournyFirstRound"        json:"tournyFirstRound"`
	TournySecondRound   bool                   `bson:"tournySecondRound"       json:"tournySecondRound"`
	Casual              bool                   `bson:"casual"                  json:"casualGame"`
	Practice            bool                   `bson:"practice"                json:"practiceGame"`
	Custom              bool                   `bson:"custom"                  json:"custom"`
	Unlisted            bool                   `bson:"unlisted"                json:"unlistedGame"`
	VerifiedOnly        bool                   `bson:"verifiedOnly"            json:"isVerifiedOnly"`
	Chats               []Chat                 `bson:"chats"                   json:"chats"`
	Guesses             map[string]string      `bson:"guesses"                 json:"guesses"`
	Timer               int                    `bson:"timer"                   json:"timedMode"`
	Blind               bool                   `bson:"blind"                   json:"blindMode"`
	Completed           bool                   `bson:"completed"               json:"completed"`
	GameState           GameState              `bson:"gameState"               json:"gameState"`
	CustomGameSettings  map[string]interface{} `bson:"customGameSettings"      json:"customGameSettings"`
	PublicPlayersState  []interface{}          `bson:"publicPlayersState"      json:"publicPlayersState"`
	GeneralGameSettings GeneralGameSettings    `bson:"general"                 json:"general"`
	PlayersState        []interface{}          `bson:"playersState"            json:"playersState"`
	CardFlingerState    []interface{}          `bson:"cardFlingerState"        json:"cardFlingerState"`
	TrackState          TrackState             `bson:"trackState"              json:"trackState"`
}

type GamePrivate struct {
	Game
	Reports                 interface{}   `bson:"reports" json:"reports"`
	UnseatedGameChats       []Chat        `bson:"unseatedGameChats" json:"unseatedGameChats"`
	CommandChats            []Chat        `bson:"commandChats" json:"commandChats"`
	ReplayGameChats         []Chat        `bson:"replayGameChats" json:"replayGameChats"`
	Lock                    interface{}   `bson:"lock" json:"lock"`
	VotesPeeked             bool          `bson:"votesPeeked" json:"votesPeeked"`
	RemakeVotesPeeked       bool          `bson:"remakeVotesPeeked" json:"remakeVotesPeeked"`
	InvIndex                int           `bson:"invIndex" json:"invIndex"`
	HiddenInfoChat          []Chat        `bson:"hiddenInfoChat" json:"hiddenInfoChat"`
	HiddenInfoSubscriptions []interface{} `bson:"hiddenInfoSubscriptions" json:"hiddenInfoSubscriptions"`
	HiddenInfoShouldNotify  bool          `bson:"hiddenInfoShouldNotify" json:"hiddenInfoShouldNotify"`
	GameCreatorName         string        `bson:"gameCreatorName" json:"gameCreatorName"`
	GameCreatorID           string        `bson:"gameCreatorID" json:"gameCreatorID"`
	GameCreatorBlacklist    []string      `bson:"gameCreatorBlacklist" json:"gameCreatorBlacklist"`
	PrivatePassword         string        `bson:"privatePassword" json:"privatePassword"`
}
