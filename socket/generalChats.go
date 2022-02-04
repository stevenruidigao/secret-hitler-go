package socket

import (
	"secrethitler.io/types"

	"sync"
)

var GeneralChats = types.GeneralChats{}
var GeneralChatsMutex = sync.RWMutex{}
