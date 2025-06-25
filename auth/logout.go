package auth

import (
	"sync"
)

var m sync.Mutex

func Logout(chatID int64) {
	m.Lock()
	defer m.Unlock()
	delete(userTokens, chatID)
	SaveTokens()
}
