package git

import "sync"

var mu sync.Mutex
var operationInProgress bool
