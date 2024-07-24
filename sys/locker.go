package sys

type Locker struct {
	c chan struct{}
}

func NewLocker() Locker {
	var l Locker
	l.c = make(chan struct{}, 1)
	l.c <- struct{}{}
	return l
}

func (l Locker) Lock() bool {
	rel := false
	select {
	case <-l.c:
		rel = true
	default:
	}
	return rel
}

func (l Locker) Unlock() {
	l.c <- struct{}{}
}
