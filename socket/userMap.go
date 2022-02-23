package socket

import (
	"secrethitler.io/types"

	"sync"
)

var UserMap = map[string]*types.UserPrivate{}
var UserMapMutex = sync.RWMutex{}
