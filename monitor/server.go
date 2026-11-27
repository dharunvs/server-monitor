package monitor

import (
	"context"
	"time"

	"github.com/go-ping/ping"
)

func PingProbe(host string) func(context.Context) bool {
	return func(ctx context.Context) bool {
		if ctx.Err()!=nil { return false }
		p,err:=ping.NewPinger(host)
		if err!=nil { return false }
		p.Count=3
		p.Timeout=3*time.Second
		// The ping library has no context API; its timeout bounds shutdown.
		if err=p.Run(); err!=nil { return false }
		return p.Statistics().PacketsRecv>0
	}
}
