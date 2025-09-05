package worker

import "sync"

type Locker interface {
	Lock(key string)
	Unlock(key string)
}

type KeyedLocker struct {
	locks sync.Map
}

func NewKeyedLocker() *KeyedLocker {
	return &KeyedLocker{}
}

func (kl *KeyedLocker) Lock(key string) {
	val, _ := kl.locks.LoadOrStore(key, &sync.Mutex{})
	lock := val.(*sync.Mutex)

	lock.Lock()
}

func (kl *KeyedLocker) Unlock(key string) {
	val, ok := kl.locks.Load(key)
	if !ok {
		return
	}

	lock := val.(*sync.Mutex)
	lock.Unlock()
}
