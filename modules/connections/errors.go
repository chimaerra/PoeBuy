package connections

import "errors"

var (
    ErrWatcherPanicked = errors.New("watcher panicked and recovered")
)