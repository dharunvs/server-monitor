package monitor

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Result struct {
	Host, Type, Parameter, Value string
}

type Check struct {
	Host, Type, Parameter string
	Interval time.Duration
	Probe func(context.Context) bool
	Alert bool
}

type Sink func(context.Context, Result) error
type Alert func(context.Context, string) error

type scheduled struct {
	check Check
	next time.Time
	busy atomic.Bool
	seen, healthy bool
}

func Run(ctx context.Context, checks []Check, workers int, sink Sink, alert Alert, report func(error)) error {
	if workers < 1 || sink == nil {
		return fmt.Errorf("positive worker count and sink required")
	}
	if report == nil { report = func(error) {} }
	tasks := make([]*scheduled, 0, len(checks))
	for _, c := range checks {
		if c.Interval <= 0 || c.Probe == nil { return fmt.Errorf("invalid check %s/%s", c.Host, c.Parameter) }
		tasks = append(tasks, &scheduled{check:c})
	}
	jobs := make(chan *scheduled, workers)
	var wg sync.WaitGroup
	for i:=0; i<workers; i++ {
		wg.Add(1)
		go func(){
			defer wg.Done()
			for {
				select {
				case <-ctx.Done(): return
				case task := <-jobs:
					if ctx.Err()!=nil { task.busy.Store(false); return }
					execute(ctx, task, sink, alert, report)
				}
			}
		}()
	}
	defer wg.Wait()
	ticker := time.NewTicker(100*time.Millisecond)
	defer ticker.Stop()
	cursor := 0
	for {
		now := time.Now()
		for offset:=0; offset<len(tasks); offset++ {
			idx := (cursor+offset)%len(tasks)
			t := tasks[idx]
			if now.Before(t.next) || !t.busy.CompareAndSwap(false,true) { continue }
			select {
			case jobs <- t:
				t.next=now.Add(t.check.Interval)
			default:
				t.busy.Store(false)
			}
		}
		if len(tasks)>0 { cursor=(cursor+1)%len(tasks) }
		select {
		case <-ctx.Done(): return nil
		case <-ticker.C:
		}
	}
}

func execute(ctx context.Context, t *scheduled, sink Sink, alert Alert, report func(error)) {
	defer t.busy.Store(false)
	c:=t.check
	healthy:=c.Probe(ctx)
	if ctx.Err()!=nil { return }
	value:="0"
	if healthy { value="1" }
	if c.Type=="Service" {
		value="Not Running"
		if healthy { value="Running" }
	}
	writeCtx,cancel:=context.WithTimeout(ctx,5*time.Second)
	err:=sink(writeCtx,Result{c.Host,c.Type,c.Parameter,value})
	cancel()
	if err!=nil { report(fmt.Errorf("write monitoring result: %w",err)) }
	if c.Alert && alert!=nil && ((!t.seen&&!healthy)||(t.seen&&t.healthy!=healthy)) {
		state:="DOWN"
		if healthy { state="RECOVERED" }
		alertCtx,cancel:=context.WithTimeout(ctx,5*time.Second)
		err:=alert(alertCtx,fmt.Sprintf("%s: %s %s %s",state,c.Host,c.Type,c.Parameter))
		cancel()
		if err!=nil { report(fmt.Errorf("send alert: %w",err)) } else { t.seen=true; t.healthy=healthy }
	} else {
		t.seen=true; t.healthy=healthy
	}
}
