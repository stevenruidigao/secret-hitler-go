package socket

import (
	"secrethitler.io/types"
	"secrethitler.io/utils"

	"fmt"
	"strconv"
	"time"
)

func StartElection(game *types.GamePrivate, specialElectionPresidentIndex int) {
	if game.GamePublic.TrackState.FascistPolicyCount >= game.GamePublic.CustomGameSettings.VetoZone {
		game.GamePublic.GameState.VetoEnabled = true
	}

	if game.GamePublic.GameState.UndrawnPolicyCount < 3 {
		ShufflePolicies(game)
	}

	game.GamePublic.GameState.PresidentIndex++

	if specialElectionPresidentIndex != -1 {
		game.GamePublic.GameState.PresidentIndex = specialElectionPresidentIndex

	} else if game.GamePublic.GameState.SpecialElectionFormerPresidentIndex != -1 {
		game.GamePublic.GameState.PresidentIndex = game.GamePublic.GameState.SpecialElectionFormerPresidentIndex + 1
		game.GamePublic.GameState.SpecialElectionFormerPresidentIndex = -1
	}

	game.GamePublic.GameState.PresidentIndex %= game.GamePublic.PlayerCount

	for game.GamePublic.PublicPlayerStates[game.GamePublic.GameState.PresidentIndex].Dead {
		game.GamePublic.GameState.PresidentIndex++
		game.GamePublic.GameState.PresidentIndex %= game.GamePublic.PlayerCount
	}

	game.GamePublic.GeneralGameSettings.ElectionCount++

	IO.BroadcastToRoom("/", "aem", "gameList", GetGameList(true))
	IO.BroadcastToRoom("/", "users", "gameList", GetGameList(false))

	game.GamePublic.GeneralGameSettings.Status = "Election #" + strconv.Itoa(game.GamePublic.GeneralGameSettings.ElectionCount) + ": president to select a chancellor"

	if !game.GamePublic.GeneralGameSettings.Experienced && !game.GamePublic.GeneralGameSettings.DisableGamechat {
		game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].GameChats = append(game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].GameChats, types.PlayerChat{
			Timestamp: time.Now(),
			GameChat:  true,
			Chat: []types.GameChat{
				types.GameChat{
					Text: "You are president and must select a chancellor.",
				},
			},
		})
	}

	var clickActionInfo []int

	for i := range game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].PlayerStates {
		if !game.SeatedPlayers[i].Dead && i != game.GamePublic.GameState.PresidentIndex && (game.GamePublic.GeneralGameSettings.LivingPlayerCount > 5 || i != game.GamePublic.GameState.PreviousElectedGovernment[0]) && i != game.GamePublic.GameState.PreviousElectedGovernment[1] {
			game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].PlayerStates[i].NotificationStatus = "notification"
			clickActionInfo = append(clickActionInfo, i)
		}
	}

	game.GamePublic.GameState.ClickActionInfo = []interface{}{game.SeatedPlayers[game.GamePublic.GameState.PresidentIndex].UserPublic.Username, clickActionInfo}
	fmt.Println("ClickActionInfo", game.GamePublic.GameState.ClickActionInfo)

	for i := range game.GamePublic.PublicPlayerStates {
		game.GamePublic.PublicPlayerStates[i].CardStatus.CardDisplayed = false
		game.GamePublic.PublicPlayerStates[i].GovernmentStatus = ""
	}

	game.GamePublic.PublicPlayerStates[game.GamePublic.GameState.PresidentIndex].GovernmentStatus = "isPendingPresident"
	game.GamePublic.PublicPlayerStates[game.GamePublic.GameState.PresidentIndex].Loader = true
	game.GamePublic.GameState.Phase = "selectingChancellor"

	if game.GamePublic.GeneralGameSettings.Timer > 0 {
		if game.Timer != nil {
			game.Timer.Stop()
			game.Timer = nil
		}

		game.GamePublic.GameState.TimedMode = true
		game.Timer = time.AfterFunc(time.Duration(game.GamePublic.GeneralGameSettings.Timer)*time.Second, func() {
			if game.GamePublic.GameState.TimedMode {
				clickActionInfo, ok := game.GamePublic.GameState.ClickActionInfo[1].([]int)

				if !ok {
					return
				}

				chancellorIndex := clickActionInfo[utils.RandInt(0, uint32(len(clickActionInfo)))]
				SelectChancellor(&game.GamePublic.PublicPlayerStates[game.GamePublic.GameState.PresidentIndex].UserPublic, game, chancellorIndex, false)
			}
		})
	}

	SendInProgressGameUpdate(game)
}
