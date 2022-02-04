package types

import (
	"time"
)

type GameFilters struct {
	Public     bool `bson:"public"     json:"public"`
	Private    bool `bson:"private"    json:"private"`
	Unstarted  bool `bson:"unstarted"  json:"unstarted"`
	InProgress bool `bson:"inProgress" json:"inProgress"`
	Completed  bool `bson:"completed"  json:"completed"`
	CustomGame bool `bson:"customGame" json:"customGame"`
	CasualGame bool `bson:"casualGame" json:"casualGame"`
	TimedMode  bool `bson:"timedMode"  json:"timedMode"`
	Standard   bool `bson:"standard"   json:"standard"`
	Rainbow    bool `bson:"rainbow"    json:"rainbow"`
}

type GameNotes struct {
	Top    int `bson:"top"    json:"top"`
	Left   int `bson:"left"   json:"left"`
	Width  int `bson:"width"  json:"width"`
	Height int `bson:"height" json:"height"`
}

type GameSettings struct {
	DisableVisibleElo       bool        `bson:"disableVisibleElo"       json:"staffDisableVisibleElo"`
	DisableStaffColor       bool        `bson:"disableStaffColor"       json:"staffDisableStaffColor"`
	Incognito               bool        `bson:"incognito"               json:"staffIncognito"`
	Rainbow                 bool        `bson:"rainbow"                 json:"isRainbow"`
	NewReport               bool        `bson:"newReport"               json:"newReport"`
	Cardback                string      `bson:"cardback"                json:"customCardback,omitempty"`
	CardbackSaveTime        string      `bson:"cardbackSaveTime"        json:"customCardbackSaveTime,omitempty"`
	CardbackID              string      `bson:"cardbackID"              json:"customCardbackUid,omitempty"`
	Timestamps              bool        `bson:"timestamps"              json:"enableTimestamps"`
	RightSidebarInGame      bool        `bson:"rightSideBarInGame"      json:"enableRightSidebarInGame"`
	PlayerColorsInChat      bool        `bson:"playerColorsInChat"      json:"playerColorsInChat"`
	PlayerCardbacks         bool        `bson:"playerCardbacks"         json:"playerCardbacks"`
	HelpMessages            bool        `bson:"helpMessages"            json:"helpMessages"`
	HelpIcons               bool        `bson:"helpIcons"               json:"helpIcons"`
	Confetti                bool        `bson:"confetti"                json:"confetti"`
	Crowns                  bool        `bson:"crowns"                  json:"crowns"`
	Seasonal                bool        `bson:"seasonal"                json:"seasonal"`
	Aggregations            bool        `bson:"aggregations"            json:"aggregations"`
	KillConfirmation        bool        `bson:"killConfirmation"        json:"killConfirmation"`
	SoundStatus             string      `bson:"soundStatus"             json:"soundStatus"`
	UnbanTime               time.Time   `bson:"unbanTime"               json:"unbanTime"`
	UnTimeoutTime           time.Time   `bson:"unTimeoutTime"           json:"unTimeoutTime"`
	FontSize                int         `bson:"fontSize"                json:"fontSize"`
	FontFamily              string      `bson:"fontFamily"              json:"fontFamily"`
	Private                 bool        `bson:"private"                 json:"isPrivate"`
	PrivateToggleTime       time.Time   `bson:"privateToggleTime"       json:"privateToggleTime"`
	Blacklist               []string    `bson:"blacklist"               json:"blacklist"`
	TournyWins              []string    `bson:"tournyWins"              json:"tournyWins"`
	ChangedName             bool        `bson:"changedName"             json:"changedName"`
	SeasonAwards            []string    `bson:"seasonAwards"            json:"seasonAwards"`
	SpecialTournamentStatus string      `bson:"specialTournamentStatus" json:"specialTournamentStatus"`
	DisableElo              bool        `bson:"disableElo"              json:"disableElo"`
	FullHeight              bool        `bson:"fullHeight"              json:"fullHeight"`
	SafeForWork             bool        `bson:"safeForWork"             json:"safeForWork"`
	KeyboardShortcuts       string      `bson:"keyboardShortcuts"       json:"keyboardShortcuts"`
	GameFilters             GameFilters `bson:"gameFilters"             json:"gameFilters"`
	GameNotes               GameNotes   `bson:"gameNotes"               json:"gameNotes"`
	PlayerNotes             []string    `bson:"playerNotes"             json:"playerNotes"`
	IgnoreIPBans            bool        `bson:"ignoreIPBans"            json:"ignoreIPBans"`
	TruncatedSize           int         `bson:"truncatedSize"           json:"truncatedSize"`
	ClaimCharacters         string      `bson:"claimCharacters"         json:"claimCharacters"`
	ClaimButtons            string      `bson:"claimButtons"            json:"claimButtons"`
	CustomWidth             string      `bson:"customWidth"             json:"customWidth"`
}
