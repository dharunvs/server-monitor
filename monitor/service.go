package monitor

import (
	"context"
	"os/exec"
	"time"
)

// ServiceProbe only observes service health. It never restarts services.
func ServiceProbe(name string) func(context.Context) bool {
	return func(ctx context.Context) bool {
		ctx,cancel:=context.WithTimeout(ctx,5*time.Second)
		defer cancel()
		return exec.CommandContext(ctx,"systemctl","is-active","--quiet","--",name).Run()==nil
	}
}
