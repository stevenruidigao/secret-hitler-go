package socket

import (
	"secrethitler.io/types"

	"sync"
)

var GameMap = map[string]*types.GamePrivate{}
var GameMapMutex = sync.RWMutex{}
