package types

import (
	"html/template"
)

type RenderData struct {
	Game                       interface{}
	ProdCacheBustToken         string
	Username                   string
	Home                       interface{}
	Changelog                  interface{}
	Rules                      interface{}
	Howtoplay                  interface{}
	Stats                      interface{}
	Wiki                       interface{}
	Discord                    interface{}
	Github                     interface{}
	Tou                        interface{}
	About                      interface{}
	PrimaryColor               template.CSS
	SecondaryColor             template.CSS
	TertiaryColor              template.CSS
	BackgroundColor            template.CSS
	SecondaryBackgroundColor   template.CSS
	TertiaryBackgroundColor    template.CSS
	TextColor                  template.CSS
	SecondaryTextColor         template.CSS
	TertiaryTextColor          template.CSS
	GameSettings               interface{}
	Verified                   interface{}
	StaffRole                  string
	HasNotDismissedSignupModal interface{}
	IsTournamentMod            interface{}
	Blacklist                  interface{}
}
