package status

import "sync"

var mapStorageStatus map[string]string // 日志仓是否正在使用
var mu sync.Mutex

func init() {
	mapStorageStatus = make(map[string]string)
}

func UpdateStorageStatus(name string, open bool) {
	mu.Lock()
	defer mu.Unlock()
	if open {
		mapStorageStatus[name] = "1"
	} else {
		delete(mapStorageStatus, name)
	}
}

func IsStorageOpening(name string) bool {
	mu.Lock()
	defer mu.Unlock()
	return mapStorageStatus[name] == "1"
}
