package socket

import (
	"sync"
)

var EmotesList = map[string]string{}
var EmotesListMutex = sync.RWMutex{}
