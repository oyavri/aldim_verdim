package worker

import "sync"

type Locker interface {
	Lock(key string)
	Unlock(key string)
}

type KeyedLocker struct {
	mu    sync.Mutex
	locks map[string]*sync.Mutex
}

func NewKeyedLocker() *KeyedLocker {
	return &KeyedLocker{
		locks: make(map[string]*sync.Mutex),
	}
}

func (kl *KeyedLocker) Lock(key string) {
	kl.mu.Lock()
	defer kl.mu.Unlock()

	lock, ok := kl.locks[key]
	if !ok {
		lock = &sync.Mutex{}
		kl.locks[key] = lock
	}
	lock.Lock()
}

func (kl *KeyedLocker) Unlock(key string) {
	kl.mu.Lock()
	defer kl.mu.Unlock()

	if lock, ok := kl.locks[key]; ok {
		lock.Unlock()
	}
}
