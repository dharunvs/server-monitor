package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"root/config"
	"root/monitor"
	"root/notifier"
)

func main() {
	if err:=run(); err!=nil { log.Print(err); os.Exit(1) }
}

func run() error {
	cfg,err:=config.LoadConfig("config.json")
	if err!=nil { return fmt.Errorf("load configuration: %w",err) }
	ctx,stop:=signal.NotifyContext(context.Background(),os.Interrupt,syscall.SIGTERM)
	defer stop()
	if cfg.Interval.Availability<=0||cfg.Interval.Port<=0||cfg.Interval.Service<=0 { return fmt.Errorf("monitoring intervals must be positive") }
	dbURL:=url.URL{Scheme:"postgres",User:url.UserPassword(cfg.Database.User,cfg.Database.Password),Host:net.JoinHostPort(cfg.Database.Host,strconv.Itoa(cfg.Database.Port)),Path:"/"+cfg.Database.Database}
	// lib/pq defaults to requiring TLS. Local development may explicitly opt out.
	if mode:=os.Getenv("MONITOR_DB_SSLMODE"); mode!="" { q:=dbURL.Query(); q.Set("sslmode",mode); dbURL.RawQuery=q.Encode() }
	db,err:=sql.Open("postgres",dbURL.String())
	if err!=nil { return fmt.Errorf("invalid database configuration") }
	defer db.Close()
	db.SetMaxOpenConns(8)
	db.SetMaxIdleConns(4)
	startup,cancel:=context.WithTimeout(ctx,5*time.Second)
	err=db.PingContext(startup)
	cancel()
	if err!=nil { return fmt.Errorf("database connection failed; verify settings and TLS") }
	checks:=[]monitor.Check{}
	addHost:=func(host string,ports []int){
		checks=append(checks,monitor.Check{Host:host,Type:"Server",Parameter:"Availability",Interval:time.Duration(cfg.Interval.Availability)*time.Second,Probe:monitor.PingProbe(host),Alert:true})
		for _,port:=range ports {
			checks=append(checks,monitor.Check{Host:host,Type:"Port",Parameter:strconv.Itoa(port),Interval:time.Duration(cfg.Interval.Port)*time.Second,Probe:monitor.TCPProbe(host,port),Alert:true})
		}
	}
	for _,server:=range cfg.Monitoring.Servers { addHost(server.ServerIp,server.Ports) }
	if cfg.SelfMonitoring.Enabled {
		addHost("localhost",cfg.SelfMonitoring.Ports)
		for _,name:=range cfg.SelfMonitoring.Services {
			checks=append(checks,monitor.Check{Host:"localhost",Type:"Service",Parameter:name,Interval:time.Duration(cfg.Interval.Service)*time.Second,Probe:monitor.ServiceProbe(name),Alert:true})
		}
	}
	sink:=func(ctx context.Context,r monitor.Result) error {
		_,err:=db.ExecContext(ctx,"INSERT INTO monitoring_data (host, type, parameter, value) VALUES ($1,$2,$3,$4)",r.Host,r.Type,r.Parameter,r.Value)
		return err
	}
	var alert monitor.Alert
	if cfg.Notifier.Telegram.Token!=""&&cfg.Notifier.Telegram.ChatID!="" {
		alert=func(ctx context.Context,message string) error { return notifier.SendContext(ctx,message,&cfg.Notifier.Telegram) }
	}
	return monitor.Run(ctx,checks,8,sink,alert,func(err error){log.Print(err)})
}
