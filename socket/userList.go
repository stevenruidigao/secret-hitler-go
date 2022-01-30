package socket

import (
	"secrethitler.io/types"

	"sync"
)

var UserList = []types.User{}
var UserListMutex = sync.RWMutex{}
