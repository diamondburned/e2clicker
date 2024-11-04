package fxhooking

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"
)

// WrapRun converts a blocking function into an fx.Hook.
func WrapRun(f func(ctx context.Context) error) fx.Hook {
	pctx, pcancel := context.WithCancelCause(context.Background())
	done := make(chan error, 1)
	return fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				defer pcancel(nil)
				done <- f(pctx)
			}()

			select {
			case err := <-done:
				if err != nil {
					return err
				}
				return fmt.Errorf("hook function returned nil error")
			case <-time.After(5 * time.Second):
				return nil
			}
		},
		OnStop: func(ctx context.Context) error {
			pcancel(ctx.Err())
			return <-done
		},
	}
}
