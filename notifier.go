package notifier

import "sync"

// Notifier is similar to a sync.Cond. It allows waiters to block on a channel
// until Notify is called so that it can be used with select.
type Notifier struct {
	m  sync.Mutex
	ch chan struct{}
}

// NewNotifier returns a new Notifier.
func NewNotifier() *Notifier {
	return &Notifier{
		ch: make(chan struct{}),
	}
}

// Notify notifies all already waiting waiters.
func (n *Notifier) Notify() {
	n.m.Lock()
	defer n.m.Unlock()
	ch := n.ch
	n.ch = make(chan struct{})
	close(ch)
}

// Wait returns a channel that will be closed when Notify is called.
func (n *Notifier) Wait() <-chan struct{} {
	n.m.Lock()
	defer n.m.Unlock()
	return n.ch
}
