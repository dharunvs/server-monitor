package monitor

import (
	"context"
	"net"
	"strconv"
	"time"
)

func TCPProbe(host string, port int) func(context.Context) bool {
	return func(ctx context.Context) bool {
		d:=net.Dialer{Timeout:5*time.Second}
		conn,err:=d.DialContext(ctx,"tcp",net.JoinHostPort(host,strconv.Itoa(port)))
		if err!=nil { return false }
		conn.Close()
		return true
	}
}
