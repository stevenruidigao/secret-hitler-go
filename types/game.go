package types

import (
	"sync"
	"time"

	"github.com/googollee/go-socket.io"
)

type PlayerChat struct {
	// GameChats
	// Message   string     `bson:"message"   json:"chat"`
	UserID                string      `bson:"userID"    json:"userID"`
	Username              string      `bson:"username"  json:"userName"`
	StaffRole             string      `bson:"staffRole" json:"staffRole"`
	Timestamp             time.Time   `bson:"timestamp" json:"timestamp"`
	GameID                string      `bson:"gameID"    json:"uid"`
	GameChat              bool        `bson:"gameChat"  json:"gameChat"`
	Chat                  interface{} `bson:"chat" json:"chat"` // string or []GameChat
	RemainingPoliciesChat bool        `bson:"remainingPoliciesChat" json:"isRemainingPolicies"`
	Seat                  int         `bson:"seat" json:"seat"`
}

type GameChat struct {
	Text     string   `bson:"text" json:"text"`
	Type     string   `bson:"type" json:"type"`
	Policies []string `bson:"policies" json:"policies"`
}

// type GameChats struct {
// 	Timestamp time.Time  `bson:"timestamp" json:"timestamp"`
// 	GameChat  bool       `bson:"gameChat" json:"gameChat"`
// 	Chat      []GameChat `bson:"chat" json:"chat"`
// }

// type CustomDeckState struct {
// 	Liberal int `bson:"lib" json:"lib"`
// 	Fascist int `bson:"fas" json:"fas"`
// }

type CustomGameCounter struct {
	Liberal int `bson:"lib" json:"lib"`
	Fascist int `bson:"fas" json:"fas"`
}

type CustomGameSettings struct {
	Enabled             bool              `bson:"enabled" json:"enabled"`
	Powers              []*string         `bson:"powers"  json:"powers"`
	HitlerZone          int               `bson:"hitlerZone" json:"hitlerZone"`
	VetoZone            int               `bson:"vetoZone" json:"vetoZone"`
	TrackState          CustomGameCounter `bson:"trackState" json:"trackState"`
	DeckState           CustomGameCounter `bson:"deckState" json:"deckState"`
	FascistCount        int               `bson:"fascistCount" json:"fascistCount"`
	LiberalCount        int               `bson:"liberalCount" json:"liberalCount"`
	HitlerKnowsFascists bool              `bson:"hitlerKnowsFascists" json:"hitKnowsFas"`
}

type GameState struct {
	PreviousElectedGovernment           []int         `bson:"previousElectedGovernment" json:"previousElectedGovernment"`
	UndrawnPolicyCount                  int           `bson:"undrawnPolicyCount"        json:"undrawnPolicyCount"`
	DiscardedPolicyCount                int           `bson:"discardedPolicyCount"      json:"discardedPolicyCount"`
	PresidentIndex                      int           `bson:"presidentIndex"            json:"presidentIndex"`
	Started                             bool          `bson:"started" json:"isStarted"`
	TracksFlipped                       bool          `bson:"tracksFlipped" json:"isTracksFlipped"`
	Veto                                bool          `bson:"veto" json:"isVetoEnabled"`
	SpecialElectionFormerPresidentIndex int           `bson:"specialElectionFormerPresidentIndex" json:"specialElectionFormerPresidentIndex"`
	Phase                               string        `bson:"phase" json:"phase"`
	TimedMode                           bool          `bson:"timedMode" json:"timedModeEnabled"`
	ClickActionInfo                     []interface{} `bson:"clickActionInfo" json:"clickActionInfo"` // Username, []int
	PendingChancellorIndex              int           `bson:"pendingChancellorIndex" json:"pendingChancellorIndex"`
	ChancellorIndex                     int           `bson:"chancellorIndex" json:"chancellorIndex"`
	Frozen                              bool          `bson:"frozen" json:"isGameFrozen"`
	AudioCue                            string        `bson:"audioCue" json:"audioCue"`
	Completed                           bool          `bson:"completed" json:"isCompleted"`
	TimeCompleted                       time.Time     `bson:"timeCompleted" json:"timeCompleted"`
}

/*cardStatus: {cardDisplayed: false, isFlipped: false, cardFront: "ballot", cardBack: {cardName: "nein"}}
connected: true
customCardBack: "png"
customCardBackUid: "jayn2v1od5q"
governmentStatus: ""
isLoader: false
isRemakeVoting: false
leftGame: false
pingTime: 1643945968329
previousGovernmentStatus: ""
userName: "evanator5000"*/

// type VoteCardStatus struct {
// 	CardDisplayed bool   `bson:"cardDisplayed" json:"cardDisplayed"`
// 	Flipped       bool   `bson:"flipped"     json:"isFlipped"`
// 	CardFront     string `bson:"cardFront"     json:"cardFront"`
// 	CardBack      struct {
// 		CardName string `bson:"cardName" json:"cardName"`
// 	} `bson:"cardBack" json:"cardBack"`
// }

type CardFlinger struct {
	Position           string     `bson:"position" json:"position"`
	NotificationStatus string     `bson:"notificationStatus" json:"notificationStatus"`
	Action             string     `bson:"action" json:"action"`
	CardStatus         CardStatus `bson:"cardStatus" json:"cardStatus"`
	Discard            bool       `bson:"discard" json:"discard"`
}

type VoteStatus struct {
	Voted    bool `bson:"voted" json:"hasVoted"`
	VotedYes bool `bson:"votedYes" json:"didVoteYes"`
}

type PlayerState struct {
	UserPublic
	CardStatus               CardStatus    `bson:"cardStatus" json:"cardStatus"`
	Connected                bool          `bson:"connected" json:"connected"`
	LeftGame                 bool          `bson:"leftGame"  json:"leftGame"`
	CustomCardBack           string        `bson:"customCardBack" json:"customCardback"`
	CustomCardBackID         string        `bson:"customCardBackID" json:"customCardbackUid"`
	Loader                   bool          `bson:"loader" json:"isLoader"`
	RemakeVoting             bool          `bson:"remakeVoting" json:"isRemakeVoting"`
	PingTime                 int64         `bson:"pingTime" json:"pingTime"`
	PreviousGovernmentStatus string        `bson:"previousGovernmentStatus" json:"previousGovernmentStatus"`
	GovernmentStatus         string        `bson:"governmentStatus"         json:"governmentStatus"`
	GameChats                []PlayerChat  `bson:"gameChats"                json:"gameChats"`
	Investigated             bool          `bson:"investigated"          json:"wasInvestigated"`
	Role                     CardBack      `bson:"role"                     json:"role"`
	PlayerStates             []PlayerState `bson:"playerStates" json:"playersState"`
	NotificationStatus       string        `bson:"notificationStatus" json:"notificationStatus"`
	NameStatus               string        `bson:"nameStatus" json:"nameStatus"`
	Index                    int           `bson:"index" json:"index"`
	Dead                     bool          `bson:"dead" json:"isDead"`
	CardFlingerState         []CardFlinger `bson:"cardFlingerState" json:"cardFlingerState"`
	VoteStatus               VoteStatus    `bson:"voteStatus" json:"voteStatus"`
	Claim                    string        `bson:"claim" json:"claim"`
	// Socket                   socketio.Conn    `bson:"-" json:"-"`
	Confetti bool `bson:"confetti" json:"isConfetti"`
	Won      bool `bson:"won" json:"wonGame"`
}

type Policy struct {
	CardBack string `bson:"cardBack" json:"cardBack"`
	Flipped  bool   `bson:"flipped" json:"isFlipped"`
	Position string `bson:"position" json:"position"`
}

type TrackState struct {
	LiberalPolicyCount   int      `bson:"liberalPolicyCount"   json:"liberalPolicyCount"`
	FascistPolicyCount   int      `bson:"fascistPolicyCount"   json:"fascistPolicyCount"`
	ElectionTrackerCount int      `bson:"electionTrackerCount" json:"electionTrackerCount"`
	EnactedPolicies      []Policy `bson:"enactedPolicies"      json:"enactedPolicies"`
}

type CardStatus struct {
	CardDisplayed bool        `bson:"cardDisplayed" json:"cardDisplayed"`
	Flipped       bool        `bson:"flipped"     json:"isFlipped"`
	CardFront     string      `bson:"cardFront"     json:"cardFront"`
	CardBack      interface{} `bson:"cardBack" json:"cardBack"`
	// CardBack string `bson:"cardBack" json:"cardBack"`
	CardName string `bson:"cardName" json:"cardName"`
}

type CardBack struct {
	CardName string `bson:"cardName" json:"cardName"`
	Icon     *int   `bson:"icon" json:"icon,omitempty"`
	Team     string `bson:"team" json:"team,omitempty"`
}

type TeamElo struct {
	Overall  float64 `bson:"overall" json:"overall"`
	Seasonal float64 `bson:"seasonal" json:"season"`
}

type GameSummary struct {
	GameID             string              `bson:"gameID" json:"uid"`
	Time               time.Time           `bson:"time" json:"time"`
	GameSettings       GeneralGameSettings `bson:"gameSettings" json:"gameSettings"`
	CustomGameSettings CustomGameSettings  `bson:"customGameSettings" json:"customGameSettings"`
	Players            []PlayerState       `bson:"players" json:"players"`
	Logs               []interface{}       `bson:"logs" json:"logs"`
	LiberalElo         TeamElo             `bson:"liberalElo" json:"libElo"`
	FascistElo         TeamElo             `bson:"fascistElo" json:"fasElo"`
}

// type Player struct {
// 	UserPublic
// 	CardStatus CardStatus `bson:"cardStatus" json:"cardStatus"`
// 	Role       CardBack             `bson:"role" json:"role"`
// 	Socket     socketio.Conn    `bson:"-" json:"-"`
// 	Connected  bool             `bson:"connected" json:"connected"`
// 	LeftGame   bool             `bson:"leftGame" json:"leftGame"`
// }

type GeneralGameSettings struct {
	WhitelistedPlayers      []string  `bson:"whitelistedPlayers"      json:"whitelistedPlayers"`
	ID                      string    `bson:"ID"                      json:"uid"`
	Name                    string    `bson:"name"                    json:"name"`
	Flag                    string    `bson:"flag"                    json:"flag"`
	MinPlayersCount         int       `bson:"minPlayersCount"         json:"minPlayersCount"`
	ExcludedPlayerCount     []int     `bson:"excludedPlayerCount"     json:"excludedPlayerCount"`
	MaxPlayersCount         int       `bson:"maxPlayersCount"         json:"maxPlayersCount"`
	Status                  string    `bson:"status"                  json:"status"`
	Experienced             bool      `bson:"experienced"             json:"experiencedMode"`
	PlayerChats             string    `bson:"playerChats"             json:"playerChats"`
	VerifiedOnly            bool      `bson:"verifiedOnly"            json:"isVerifiedOnly"`
	DisableObserverLobby    bool      `bson:"disableObserverLobby"    json:"disableObserverLobby"`
	DisableObserver         bool      `bson:"disableObserver"         json:"disableObserver"`
	Tourny                  bool      `bson:"tourny"                  json:"isTourny"`
	LastModPing             int       `bson:"lastModPing"             json:"lastModPing"`
	ChatReplTime            []int     `bson:"chatReplTime"            json:"chatReplTime"`
	DisableGamechat         bool      `bson:"disableGamechat"         json:"disableGamechat"`
	Rainbow                 bool      `bson:"rainbow"                 json:"rainbowgame"`
	Blind                   bool      `bson:"blind"                   json:"blindMode"`
	Timer                   int       `bson:"timer"           	       json:"timedMode"`
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
	CustomCardBack          []string  `bson:"customCardBack" json:"customCardback"`
	CustomCardBackUID       []string  `bson:"customCardBackUID" json:"customCardbackUid"`
	// Players                 []Player               `bson:"players" json:"players"`
	SeatedCount   int        `bson:"seatedCount" json:"seatedCount"`
	TimeAbandoned *time.Time `bson:"timeAbandoned" json:"timeAbandoned"`
	// Mutex         *sync.RWMutex `bson:"mutex" json:"mutex"`
	// Map             map[string]interface{} `bson:"map" json:"map"`
	TimeStarted       time.Time `bson:"timeStarted" json:"timeStarted"`
	GameCreatorName   string    `bson:"gameCreatorName" json:"gameCreatorName,omitempty"`
	GameStatus        string    `bson:"gameStatus" json:"gameStatus"`
	LivingPlayerCount int       `bson:"livingPlayerCount" json:"livingPlayerCount"`
	TournyInfo        struct {
		Cancelled bool `bson:"cancelled" json:"isCancelled"`
		Round     int  `bson:"round" json:"round"`
	} `bson:"tournyInfo" json:"tournyInfo"`
	Recorded    bool `bson:"recorded" json:"isRecorded"`
	PlayerCount int  `bson:"playerCount" json:"playerCount"` // dupe
}

type GamePublic struct {
	// ID                string            `bson:"id"                      json:"uid"`
	// Name              string            `bson:"name"                    json:"name"`
	// Flag              string            `bson:"flag"                    json:"flag"`
	Date             time.Time `bson:"date"                    json:"date"`
	PlayerChats      string    `bson:"playerChats"             json:"playerChats"`
	PlayerCount      int       `bson:"playerCount"             json:"playerCount"`
	WinningPlayerIDs []string  `bson:"winningPlayerIDs"          json:"winningPlayers"`
	LosingPlayerIDs  []string  `bson:"losingPlayerIDs"           json:"losingPlayers"`
	WinningTeam      string    `bson:"winningTeam"             json:"winningTeam"`
	Season           int       `bson:"season"                  json:"season"`
	// Rainbow           bool              `bson:"rainbow"                 json:"rainbowgame"`
	// EloMinimum        int               `bson:"eloMinimum"              json:"eloMinimum"`
	// EloMaximum        int               `bson:"eloMaximum"              json:"eloMaximum"`
	// Rebalanced        string            `bson:"rebalanced"              json:"rebalanced"`
	// Rebalance6p       bool              `bson:"rebalance6p"             json:"rebalance6p"`
	// Rebalance7p       bool              `bson:"rebalance7p"             json:"rebalance7p"`
	// Rebalance9p       bool              `bson:"rebalance9p"             json:"rebalance9p"`
	// Rebalance9p2f     bool              `bson:"rebalance9p2f"           json:"rebalance9p2f"`
	// TournyFirstRound  bool              `bson:"tournyFirstRound"        json:"tournyFirstRound"`
	// TournySecondRound bool              `bson:"tournySecondRound"       json:"tournySecondRound"`
	// Casual            bool              `bson:"casual"                  json:"casualGame"`
	// Practice          bool              `bson:"practice"                json:"practiceGame"`
	// Custom            bool              `bson:"custom"                  json:"custom"`
	// Unlisted          bool              `bson:"unlisted"                json:"unlistedGame"`
	// VerifiedOnly      bool              `bson:"verifiedOnly"            json:"isVerifiedOnly"`
	Chats   []PlayerChat      `bson:"chats"                   json:"chats"`
	Guesses map[string]string `bson:"guesses"                 json:"guesses"`
	// Timer               int                 `bson:"timer"                   json:"timedMode"`
	// Blind              bool               `bson:"blind"                   json:"blindMode"`
	Completed           bool                `bson:"completed"               json:"completed"`
	GameState           GameState           `bson:"gameState"               json:"gameState"`
	CustomGameSettings  CustomGameSettings  `bson:"customGameSettings"      json:"customGameSettings"`
	PublicPlayerStates  []PlayerState       `bson:"publicPlayerStates"      json:"publicPlayersState"`
	GeneralGameSettings GeneralGameSettings `bson:"general"                 json:"general"`
	// PlayerStates        []PlayerState       `bson:"playerStates"            json:"playersState"`
	CardFlingerState        []CardFlinger  `bson:"cardFlingerState"        json:"cardFlingerState"`
	TrackState              TrackState     `bson:"trackState"              json:"trackState"`
	PlayerCounts            []int          `bson:"playerCounts" json:"playerCounts"`
	PlayerMap               map[string]int `bson:"playerMap" json:"playerMap"`
	PlayerMapMutex          *sync.RWMutex  `bson:"-" json:"-"`
	ChatMutex               *sync.RWMutex  `bson:"-" json:"-"`
	PublicPlayerStatesMutex *sync.RWMutex  `bson:"-" json:"-"`
}

type GamePrivate struct {
	GamePublic        `bson:"gamePublic"`
	Reports           []Report     `bson:"reports" json:"reports"`
	UnseatedGameChats []PlayerChat `bson:"unseatedGameChats" json:"unseatedGameChats"`
	CommandChats      []PlayerChat `bson:"commandChats" json:"commandChats"`
	ReplayGameChats   []PlayerChat `bson:"replayGameChats" json:"replayGameChats"`
	Lock              struct {
		SelectPresidentPolicy                   bool `bson:"selectPresidentPolicy" json:"selectPresidentPolicy"`
		SelectChancellorPolicy                  bool `bson:"selectChancellorPolicy" json:"selectChancellorPolicy"`
		SelectChancellor                        bool `bson:"selectChancellor" json:"selectChancellor"`
		SelectChancellorVoteOnVeto              bool `bson:"selectChancellorVoteOnVeto" json:"selectChancellorVoteOnVeto"`
		PolicyPeek                              bool `bson:"policyPeek" json:"policyPeek"`
		PolicyPeekAndDrop                       bool `bson:"policyPeekAndDrop" json:"policyPeekAndDrop"`
		SelectPlayerToExecute                   bool `bson:"selectPlayerToExecute" json:"selectPlayerToExecute"`
		ExecutePlayer                           bool `bson:"executePlayer" json:"executePlayer"`
		SelectSpecialElection                   bool `bson:"selectSpecialElection" json:"selectSpecialElection"`
		SpecialElection                         bool `bson:"specialElection" json:"specialElection"`
		SelectPartyMembershipInvestigate        bool `bson:"selectPartyMembershipInvestigate" json:"selectPartyMembershipInvestigate"`
		InvestigateLoyalty                      bool `bson:"investigateLoyalty" json:"investigateLoyalty"`
		ShowPlayerLoyalty                       bool `bson:"showPlayerLoyalty" json:"showPlayerLoyalty"`
		SelectPartyMembershipInvestigateReverse bool `bson:"selectPartyMembershipInvestigateReverse" json:"selectPartyMembershipInvestigateReverse"`
		SelectPolicies                          bool `bson:"selectPolicies"	json:"selectPolicies"`
		SelectOnePolicy                         bool `bson:"selectOnePolicy" json:"selectOnePolicy"`
		SelectBurnCard                          bool `bson:"selectBurnCard" json:"selectBurnCard"`
	} `bson:"-" json:"-"`
	VotesPeeked              bool                       `bson:"votesPeeked" json:"votesPeeked"`
	RemakeVotesPeeked        bool                       `bson:"remakeVotesPeeked" json:"remakeVotesPeeked"`
	InvIndex                 int                        `bson:"invIndex" json:"invIndex"`
	HiddenInfoChat           []PlayerChat               `bson:"hiddenInfoChat" json:"hiddenInfoChat"`
	HiddenInfoSubscriptions  []interface{}              `bson:"hiddenInfoSubscriptions" json:"hiddenInfoSubscriptions"`
	HiddenInfoShouldNotify   bool                       `bson:"hiddenInfoShouldNotify" json:"hiddenInfoShouldNotify"`
	GameCreatorName          string                     `bson:"gameCreatorName" json:"gameCreatorName"`
	GameCreatorID            string                     `bson:"gameCreatorID" json:"gameCreatorID"`
	GameCreatorBlacklist     []string                   `bson:"gameCreatorBlacklist" json:"gameCreatorBlacklist"`
	PrivatePassword          string                     `bson:"privatePassword" json:"privatePassword"`
	SeatedPlayers            []PlayerState              `bson:"seatedPlayers" json:"seatedPlayers"`
	PlayerStates             []PlayerState              `bson:"playerStates"            json:"playersState"`
	SocketMap                map[string][]socketio.Conn `bson:"-" json:"-"`
	SocketMapMutex           *sync.RWMutex              `bson:"-" json:"-"`
	Policies                 []string                   `bson:"policies" json:"policies"`
	Timer                    *time.Timer                `bson:"-" json:"-"`
	ActionMutex              *sync.Mutex                `bson:"-" json:"-"`
	Summary                  GameSummary                `bson:"summary" json:"summary"`
	CurrentElectionPolicies  []string                   `bson:"currentElectionPolicies" json:"currentElectionPolicies"`
	CurrentChancellorOptions []string                   `bson:"currentChancellorOptions" json:"currentChancellorOptions"`
	UnsentReports            []Report                   `bson:"unsentReports" json:"unsentReports"`
}
