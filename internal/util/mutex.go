package util

import (
	"log"
	"time"
)

type TriableMutex struct {
	ch    chan struct{}
	await chan bool
}

func NewTryableMutex() *TriableMutex {
	return &TriableMutex{
		ch:    make(chan struct{}, 1),
		await: make(chan bool, 1),
	}
}

// checks if lock is acquired without blocking
func (m *TriableMutex) Locked() bool {
	return len(m.ch) > 0
}

// acquire lock, blocks until the lock is acquired
func (m *TriableMutex) Lock() {
	m.ch <- struct{}{}
}

// tries to acquire lock, blocks until the lock is acquired or timeout passed
// puts 'true' on returning channel if lock is acquired successfully, or else 'false'
func (m *TriableMutex) LockTimeout(timeout time.Duration) chan bool {
	m.disposeAwait()
	select {
	case m.ch <- struct{}{}:
		m.await <- true
	case <-time.After(timeout):
		m.await <- false
	}
	return m.await
}

// tries to acquire lock, blocks until the lock is acquired or deadline exceeded
// puts 'true' on returning channel if lock is acquired successfully, or else 'false'
func (m *TriableMutex) LockDeadline(deadline time.Time) chan bool {
	return m.LockTimeout(time.Until(deadline))
}

// release lock, blocks until the lock is released
func (m *TriableMutex) Unlock() {
	m.disposeAwait()
	<-m.ch
}

// tries to release lock, blocks until the lock is released or timeout passed
// puts 'true' on returning channel if lock is released successfully, or else 'false'
func (m *TriableMutex) UnlockTimeout(timeout time.Duration) chan bool {
	m.disposeAwait()
	select {
	case <-m.ch:
		m.await <- true
	case <-time.After(timeout):
		m.await <- false
	}
	return m.await
}

// tries to release lock, blocks until the lock is released or deadline exceeded
// puts 'true' on returning channel if lock is acquired successfully, or else 'false'
func (m *TriableMutex) UnlockDeadline(deadline time.Time) chan bool {
	return m.UnlockTimeout(time.Until(deadline))
}

func (m *TriableMutex) disposeAwait() {
	if len(m.await) == 0 {
		return
	}

	for await := range m.await {
		log.Println(await)
	}
}
