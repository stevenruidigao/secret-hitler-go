package types

import (
	"time"
)

type GameFilters struct {
	Public     bool
	Private    bool
	Unstarted  bool
	InProgress bool
	Completed  bool
	CustomGame bool
	CasualGame bool
	TimedMode  bool
	Standard   bool
	Rainbow    bool
}

type GameNotes struct {
	Top    int
	Left   int
	Width  int
	Height int
}

type GameSettings struct {
	DisableVisibleElo       bool
	DisableStaffColor       bool
	Incognito               bool
	Rainbow                 bool
	NewReport               bool
	Cardback                *string
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
	UnbanTime               time.Time
	Private                 bool
	PrivateToggleTime       time.Time
	Blacklist               []string
	TournyWins              []string
	ChangedName             bool
	SeasonAwards            []string
	SpecialTournamentStatus string
	DisableElo              bool
	FullHeight              bool
	SafeForWork             bool
	KeyboardShortcuts       string
	GameFilters             GameFilters
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
