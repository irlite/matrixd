package session

import (
	"sync/atomic"

	"github.com/irlite/matrixd/internal/wire"
)

type Notifier func(connName string, out wire.Outgoing)

var notifier atomic.Pointer[Notifier]

func SetNotify(fn Notifier) {
	notifier.Store(&fn)
}

func Notify(connName string, out wire.Outgoing) {
	fn := notifier.Load()
	if fn == nil {
		return
	}
	(*fn)(connName, out)
}
