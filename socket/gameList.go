package socket

import (
	"secrethitler.io/types"

	"sync"
)

var GameList = []types.GamePrivate{}
var GameListMutex = sync.RWMutex{}
