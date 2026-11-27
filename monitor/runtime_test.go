package monitor

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestAlertTransitions(t *testing.T) {
	var healthy bool
	var alerts []string
	task:=&scheduled{check:Check{Host:"local",Type:"Server",Parameter:"Availability",Alert:true,Probe:func(context.Context)bool{return healthy}}}
	sink:=func(context.Context,Result)error{return nil}
	alert:=func(_ context.Context,s string)error{alerts=append(alerts,s);return nil}
	for i:=0;i<3;i++ { execute(context.Background(),task,sink,alert,func(err error){t.Error(err)}) }
	healthy=true
	execute(context.Background(),task,sink,alert,func(err error){t.Error(err)})
	if len(alerts)!=2||alerts[0]!="DOWN: local Server Availability"||alerts[1]!="RECOVERED: local Server Availability" { t.Fatalf("alerts = %v",alerts) }
}

func TestBoundedWorkersAndCancellation(t *testing.T) {
	ctx,cancel:=context.WithTimeout(context.Background(),350*time.Millisecond)
	defer cancel()
	var active,peak,calls atomic.Int32
	probe:=func(ctx context.Context)bool{
		n:=active.Add(1)
		defer active.Add(-1)
		calls.Add(1)
		for { old:=peak.Load(); if n<=old||peak.CompareAndSwap(old,n){break} }
		select {case <-ctx.Done(): case <-time.After(20*time.Millisecond):}
		return true
	}
	checks:=make([]Check,20)
	for i:=range checks { checks[i]=Check{Interval:time.Millisecond,Probe:probe} }
	if err:=Run(ctx,checks,3,func(context.Context,Result)error{return nil},nil,func(err error){t.Error(err)});err!=nil { t.Fatal(err) }
	if calls.Load()==0||peak.Load()>3||active.Load()!=0 { t.Fatalf("calls=%d peak=%d active=%d",calls.Load(),peak.Load(),active.Load()) }
}

func TestInvalidCheck(t *testing.T) {
	if err:=Run(context.Background(),[]Check{{}},1,func(context.Context,Result)error{return nil},nil,nil);err==nil { t.Fatal("expected validation failure") }
}

func BenchmarkExecute(b *testing.B) {
	task:=&scheduled{check:Check{Host:"local",Probe:func(context.Context)bool{return true}}}
	sink:=func(context.Context,Result)error{return nil}
	b.ReportAllocs()
	for i:=0;i<b.N;i++ { execute(context.Background(),task,sink,nil,func(error){}) }
}
