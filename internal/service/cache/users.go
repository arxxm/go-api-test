package cache

import (
	"go-api-test/pkg/domain"
	"sync"
	"time"
)

type UsersCache struct {
	users map[int64]cachedUser
	sync.RWMutex
}

type cachedUser struct {
	user      domain.User
	createdAt time.Time
}

func NewUsersCache() *UsersCache {
	uc := &UsersCache{users: map[int64]cachedUser{}}
	go uc.cleaner()
	return uc
}

func (uc *UsersCache) Set(id int64, user domain.User) {
	uc.Lock()
	cu := cachedUser{
		user:      user,
		createdAt: time.Now(),
	}
	uc.users[id] = cu
	uc.Unlock()
}

func (uc *UsersCache) Get(id int64) (domain.User, bool) {
	uc.RLock()
	defer uc.RUnlock()
	k, v := uc.users[id]
	return k.user, v
}

func (uc *UsersCache) Delete(id int64) {
	uc.Lock()
	delete(uc.users, id)
	uc.Unlock()
}

func (uc *UsersCache) cleaner() {
	for {
		uc.Lock()
		tn := time.Now()
		diff := time.Minute * 5
		for key, u := range uc.users {
			if tn.Sub(u.createdAt) > diff {
				delete(uc.users, key)
			}
		}
		uc.Unlock()
		<-time.After(time.Second * 5)
	}
}
