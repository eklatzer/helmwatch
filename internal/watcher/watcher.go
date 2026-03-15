package watcher

import (
	"context"
	"fmt"
	"log"

	"github.com/sgtdi/fswatcher"
)

func Watch(ctx context.Context, handleEventFunc func(fswatcher.WatchEvent), opts ...fswatcher.WatcherOpt) error {
	w, err := fswatcher.New(opts...)
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}

	go func() {
		if err := w.Watch(ctx); err != nil {
			log.Printf("failed to watch: %v", err)
		}
	}()

	for event := range w.Events() {
		handleEventFunc(event)
	}

	return nil
}
