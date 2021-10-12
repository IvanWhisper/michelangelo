package generate

import "sync"

var locker *sync.Mutex
var loker_map sync.Map

func init() {
	locker = new(sync.Mutex)
	//loker_map=make(map[string]*sync.Mutex,0)
}
