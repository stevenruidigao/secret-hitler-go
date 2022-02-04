package socket

import (
	"secrethitler.io/types"

	"sync"
)

var UserMap = map[string]types.UserPublic{}
var UserMapMutex = sync.RWMutex{}
