package routes

import (
	"secrethitler.io/database"
	"secrethitler.io/types"

	"fmt"
	"html/template"
	"net/http"
)

func Render(tmplName string) http.Handler {
	fmt.Println("./views/" + tmplName + ".tmpl")
	tmpl := template.Must(template.ParseFiles("./views/" + tmplName + ".tmpl")).Funcs(template.FuncMap{
		"marshal": Marshal,
	})

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		session := Authenticate(request)
		user := types.UserPrivate{}
		// fmt.Println(session, user, "*")
		// fmt.Println("S**********:", session)

		if session != nil {
			result := database.GetUserByID(session.UserID)

			if result != nil && result.FinishedSignup {
				user = *result
			}
		}

		// fmt.Println("UUUUUUUUUUUU", user)

		data := types.RenderData{
			ProdCacheBustToken:       CacheToken,
			Username:                 user.UserPublic.Username,
			Game:                     tmplName == "game",
			Home:                     tmplName == "page-home",
			Changelog:                tmplName == "page-changelog",
			Rules:                    tmplName == "page-rules",
			Howtoplay:                tmplName == "page-howtoplay",
			Stats:                    tmplName == "page-stats",
			Wiki:                     tmplName == "page-wiki",
			Discord:                  false,
			Github:                   false,
			Tou:                      tmplName == "page-tou",
			About:                    tmplName == "page-about",
			PrimaryColor:             template.CSS("hsl(225, 73%, 57%)"),
			SecondaryColor:           template.CSS("hsl(225, 48%, 57%)"),
			TertiaryColor:            template.CSS("hsl(265, 73%, 57%)"),
			BackgroundColor:          template.CSS("hsl(0, 0%, 0%)"),
			SecondaryBackgroundColor: template.CSS("hsl(0, 0%, 7%)"),
			TertiaryBackgroundColor:  template.CSS("hsl(0, 0%, 14%)"),
			TextColor:                template.CSS("hsl(0, 0%, 100%)"),
			SecondaryTextColor:       template.CSS("hsl(0, 0%, 93%)"),
			TertiaryTextColor:        template.CSS("hsl(0, 0%, 86%)"),
			GameSettings: types.GameSettings{
				CustomWidth: "",
				FontFamily:  "",
			},
			Verified:                   false,
			StaffRole:                  user.UserPublic.StaffRole,
			HasNotDismissedSignupModal: !user.DismissedSignupModal,
			IsTournamentMod:            false,
			Blacklist:                  struct{}{},
		}

		/*if tmplName == "game" {
			data.Game = true

		} else if tmplName == "page-home" {
			data.Home = true
		}*/

		//data.ProdCacheBustToken = CacheToken
		tmpl.Execute(writer, data)
	})
}
