package socket

import (
	"secrethitler.io/types"

	"fmt"
	"strings"
	"time"
)

func CompleteGame(game *types.GamePrivate, winningTeam string) {
	if len(game.UnsentReports) > 0 {
		for i := range game.UnsentReports {
			game.UnsentReports[i].ReportType += "delayed"
			MakeReport(game, game.UnsentReports[i])
		}

		game.UnsentReports = []types.Report{}
	}

	IO.BroadcastToRoom("/", "game-"+game.GamePublic.GeneralGameSettings.ID, "removeClaim")

	if game.GamePublic.GeneralGameSettings.Timer > 0 && game.Timer != nil {
		game.Timer.Stop()
		game.Timer = nil
		game.GamePublic.GameState.TimedMode = false
	}

	if game.GamePublic.GeneralGameSettings.Recorded {
		fmt.Println("A game attempted to be rerecorded!", game.GamePublic.GeneralGameSettings.ID)
		return
	}

	if !(game.GamePublic.GeneralGameSettings.Tourny && game.GamePublic.GeneralGameSettings.TournyInfo.Round == 1) {
		for i := range game.SeatedPlayers {
			if game.SeatedPlayers[i].Role.Team == winningTeam {
				game.GamePublic.PublicPlayerStates[i].NotificationStatus = "success"
				game.GamePublic.PublicPlayerStates[i].Confetti = true
				game.GamePublic.PublicPlayerStates[i].Won = true
			}
		}

		time.AfterFunc(15*time.Second, func() {
			for i := range game.SeatedPlayers {
				if game.SeatedPlayers[i].Role.Team == winningTeam {
					game.GamePublic.PublicPlayerStates[i].Confetti = false
				}
			}

			SendInProgressGameUpdate(game)
		})
	}

	game.GamePublic.GeneralGameSettings.Status = strings.Title(winningTeam) + "s win the game."
	game.GamePublic.GameState.Completed = true
	game.GamePublic.GameState.TimeCompleted = time.Now()
	IO.BroadcastToRoom("/", "aem", "gameList", GetGameList(true))
	IO.BroadcastToRoom("/", "users", "gameList", GetGameList(false))

	for i := range game.GamePublic.PublicPlayerStates {
		game.GamePublic.PublicPlayerStates[i].NameStatus = game.SeatedPlayers[i].Role.CardName
	}

	remainingPolicies := []string{}

	for i := range game.Policies {
		if game.Policies[i] == "fascist" {
			remainingPolicies = append(remainingPolicies, "r")

		} else {
			remainingPolicies = append(remainingPolicies, "b")
		}
	}

	gameChats := []types.PlayerChat{
		types.PlayerChat{
			Timestamp: time.Now(),
			GameChat:  true,
			Chat: []types.GameChat{
				types.GameChat{
					Text: strings.Title(winningTeam) + "s",
					Type: winningTeam,
				},
				types.GameChat{
					Text: " win the game.",
				},
			},
		},
		types.PlayerChat{
			Timestamp:             time.Now(),
			RemainingPoliciesChat: true,
			Chat: []types.GameChat{
				types.GameChat{
					Text: "The remaining policies are ",
				},
				types.GameChat{
					Policies: remainingPolicies,
				},
				types.GameChat{
					Text: ".",
				},
			},
		},
	}

	for i := range game.SeatedPlayers {
		game.SeatedPlayers[i].GameChats = append(game.SeatedPlayers[i].GameChats, gameChats...)
	}

	game.UnseatedGameChats = append(game.UnseatedGameChats, gameChats...)
	SendInProgressGameUpdate(game)
	SaveGame(game)
	game.GamePublic.GeneralGameSettings.Recorded = true
}
