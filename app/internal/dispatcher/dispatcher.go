package dispatcher

import (
	"context"

	"github.com/mzeahmed/noticoel/internal/notifier"
)

type Dispatcher struct {
	notifiers []notifier.Notifier
}

func New() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) Register(n notifier.Notifier) {
	d.notifiers = append(d.notifiers, n)
}

func (d *Dispatcher) Dispatch(
	ctx context.Context,
	msg notifier.Message,
) []notifier.Result {

	results := make([]notifier.Result, 0, len(d.notifiers))

	for _, n := range d.notifiers {
		results = append(results, n.Notify(ctx, msg))
	}

	return results
}
