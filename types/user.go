package types

import (
	"time"

	"github.com/markbates/goth"
)

type Season struct {
	Wins          int `bson:"wins"          json:"wins"`
	Losses        int `bson:"losses"	json:"losses"`
	RainbowWins   int `bson:"rainbowWins"   json:"rainbowWins"`
	RainbowLosses int `bson:"rainbowLosses" json:"rainbowLosses"`
}

type Warning struct {
	Text         string    `bson:"text"         json:"text"`
	Moderator    string    `bson:"moderator"    json:"moderator"`
	Time         time.Time `bson:"time"         json:"time"`
	Acknowledged bool      `bson:"acknowledged" json:"acknowledged"`
}

type FeedbackSubmission struct {
	Text string    `bson:"text" json:"text"`
	User string    `bson:"user" json:"user"`
	Time time.Time `bson:"time" json:"time"`
}

type UserStatus struct {
	Type   string `bson:"type" json:"type"`
	GameID string `bson:"gameId" json:"gameId"`
}

type UserPublic struct {
	UserID               string               `bson:"userID"               json:"userID"`
	Username             string               `bson:"username"             json:"userName"`
	Local                bool                 `bson:"local"                json:"local,omitempty"`
	StaffRole            string               `bson:"staffRole"            json:"staffRole"`
	Contributor          bool                 `bson:"contributor"          json:"contributor"`
	DismissedSignupModal bool                 `bson:"dismissedSignupModal" json:"dismissedSignupModal,omitempty"`
	GameSettings         GameSettings         `bson:"gameSettings"         json:"gameSettings,omitempty"`
	Verified             bool                 `bson:"verified"             json:"verified,omitempty"`
	Banned               bool                 `bson:"banned"               json:"banned,omitempty"`
	Timeout              time.Time            `bson:"timeout"              json:"timeout,omitempty"`
	TOULastAgreed        string               `bson:"TOULastAgreed"	       json:"TOULastAgreed,omitempty"`
	Bio                  string               `bson:"bio"                  json:"bio,omitempty"`
	Games                []string             `bson:"games"                json:"games,omitempty"`
	Wins                 int                  `bson:"wins"                 json:"wins"`
	Losses               int                  `bson:"losses"               json:"losses"`
	RainbowWins          int                  `bson:"rainbowWins"          json:"rainbowWins"`
	RainbowLosses        int                  `bson:"rainbowLosses"        json:"rainbowLosses"`
	Seasons              []Season             `bson:"seasons"              json:"seasons,omitempty"`
	PreviousDayElo       int                  `bson:"previousDayElo"       json:"previousDayElo,omitempty"`
	Created              time.Time            `bson:"created"              json:"created,omitempty"`
	OnFire               bool                 `bson:"onFire"               json:"onFire"`
	LastCompleteGame     time.Time            `bson:"lastCompleteGame"     json:"lastCompleteGame,omitempty"`
	LastVersionSeen      string               `bson:"lastVersionSeen"      json:"lastVersionSeen,omitempty"`
	Fixed                bool                 `bson:"fixed"                json:"fixed"`
	EloSeason            float64              `bson:"eloSeason"            json:"eloSeason,omitempty"`
	EloOverall           float64              `bson:"eloOverall"           json:"eloOverall"`
	HashID               string               `bson:"hashID"               json:"hashID"`
	FeedbackSubmissions  []FeedbackSubmission `bson:"feedbackSubmissions"  json:"feedbackSubmissions,omitempty"`
	PrimaryColor         string               `bson:"primaryColor"         json:"primaryColor,omitempty"`
	SecondaryColor       string               `bson:"secondaryColor"       json:"secondaryColor,omitempty"`
	TertiaryColor        string               `bson:"tertiaryColor"        json:"tertiaryColor,omitempty"`
	BackgroundColor      string               `bson:"backgroundColor"      json:"backgroundColor,omitempty"`
	TextColor            string               `bson:"textColor"            json:"textColor,omitempty"`
	TournamentMod        bool                 `bson:"tournamentMod"        json:"tournamentMod,omitempty"`
	Status               UserStatus           `bson:"status"               json:"status,omitempty"`
	TimeLastGameCreated  time.Time            `bson:"timeLastGameCreated"  json:"timeLastGameCreated"`
	Profile              Profile              `bson:"profile"              json:"profile,omitempty"`
	LeftGame             bool                 `bson:"leftGame" json:"leftGame"`
	Connected            bool                 `bson:"connected" json:"connected"`
}

type UserPrivate struct {
	UserPublic           `bson:"userPublic" json:"userPublic"`
	LinkedAccounts       []goth.User          `bson:"linkedAccounts"       json:"linkedAccounts,omitempty"`
	Sessions             []Session            `bson:"sessions"             json:"sessions,omitempty"`
	Password             string               `bson:"password"             json:"password,omitempty"`
	Local                bool                 `bson:"local"                json:"local,omitempty"`
	DismissedSignupModal bool                 `bson:"dismissedSignupModal" json:"dismissedSignupModal,omitempty"`
	GameSettings         GameSettings         `bson:"gameSettings"         json:"gameSettings,omitempty"`
	Email                string               `bson:"email"                json:"email,omitempty"`
	SignupIP             string               `bson:"signupIP"             json:"signupIP,omitempty"`
	LastIP               string               `bson:"lastIP"               json:"lastIP,omitempty"`
	IPHistory            []string             `bson:"IPHistory"            json:"IPHistory,omitempty"`
	Verified             bool                 `bson:"verified"             json:"verified,omitempty"`
	Banned               bool                 `bson:"banned"               json:"banned,omitempty"`
	Timeout              time.Time            `bson:"timeout"              json:"timeout,omitempty"`
	TOULastAgreed        string               `bson:"TOULastAgreed"	       json:"TOULastAgreed,omitempty"`
	Bio                  string               `bson:"bio"	               json:"bio,omitempty"`
	Games                []string             `bson:"games"	               json:"games,omitempty"`
	Seasons              []Season             `bson:"seasons"              json:"seasons,omitempty"`
	PreviousDayElo       int                  `bson:"previousDayElo"       json:"previousDayElo,omitempty"`
	OnFire               bool                 `bson:"onFire"               json:"onFire"`
	LastCompleteGame     time.Time            `bson:"lastCompleteGame"     json:"lastCompleteGame,omitempty"`
	LastVersionSeen      string               `bson:"lastVersionSeen"      json:"lastVersionSeen,omitempty"`
	Fixed                bool                 `bson:"fixed"	               json:"fixed"`
	HashID               string               `bson:"hashID"               json:"hashID"`
	Warnings             []Warning            `bson:"warnings"             json:"warnings,omitempty"`
	FeedbackSubmissions  []FeedbackSubmission `bson:"feedbackSubmissions"  json:"feedbackSubmissions,omitempty"`
	PrimaryColor         string               `bson:"primaryColor"         json:"primaryColor,omitempty"`
	SecondaryColor       string               `bson:"secondaryColor"       json:"secondaryColor,omitempty"`
	TertiaryColor        string               `bson:"tertiaryColor"        json:"tertiaryColor,omitempty"`
	BackgroundColor      string               `bson:"backgroundColor"      json:"backgroundColor,omitempty"`
	TextColor            string               `bson:"textColor"            json:"textColor,omitempty"`
	TournamentMod        bool                 `bson:"tournamentMod"        json:"tournamentMod,omitempty"`
	TimeLastGameCreated  time.Time            `bson:"timeLastGameCreated"  json:"timeLastGameCreated"`
	Profile              Profile              `bson:"profile"              json:"profile"`
}

type ExtendedUser struct {
	UserPublic
	DisableVisibleElo       bool
	DisableStaffColor       bool
	Incognito               bool
	Rainbow                 bool
	NewReport               bool
	Cardback                *string `bson:"cardback"               json:"cardback"`
	CardbackSaveTime        string
	CardbackID              string
	Timestamps              bool
	RightSidebarInGame      bool
	PlayerColorsInChat      bool
	PlayerCardbacks         bool
	HelpMessages            bool
	HelpIcons               bool
	Confetti                bool
	Crowns                  bool
	Seasonal                bool
	Aggregations            bool
	KillConfirmation        bool
	SoundStatus             string
	UnbanTime               time.Time `bson:"unbanTime" json:"unbanTime"`
	Private                 bool      `bson:"private"              json:"private,omitempty"`
	PrivateToggleTime       time.Time
	Blacklist               map[string]string
	TournyWins              []string
	ChangedName             bool
	PreviousSeasonAward     string
	SpecialTournamentStatus string
	DisableElo              bool
	FullHeight              bool
	SafeForWork             bool
	KeyboardShortcuts       string
	GameFilters             GameFilters `bson:"gameFilters" json:"gameFilters"`
	GameNotes               GameNotes
	PlayerNotes             []string
	IgnoreIPBans            bool
	TruncatedSize           int
	ClaimCharacters         string
	ClaimButtons            string
	CustomWidth             string
	FontSize                int
	FontFamily              string
}
