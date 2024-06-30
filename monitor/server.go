package monitor

import (
	"fmt"
	"time"
	"sync"
	"net"
	"os/exec"
	"github.com/go-ping/ping"

	"root/logger"
	"root/config"
)


func MonitorAvailability(cfg *config.Config, server config.Server, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		pinger, err := ping.NewPinger(server.ServerIp)
		if err != nil {
			logger.Error("Failed to ping:", server.ServerIp, err)
			time.Sleep(5 * time.Second)
			continue
		}

		packetsToSend := 3
		pinger.Count = packetsToSend
    	pinger.Timeout = time.Duration(packetsToSend) * time.Second

    	err = pinger.Run() 
    	if err != nil {
	        logger.Error("Error while pinging:", server.ServerIp, err)
    	}

		stats := pinger.Statistics()
		availability := 0

		if stats.PacketsRecv > 0 {
			availability = 1
		} 
		logger.Info("Availability for ip:", server.ServerIp, availability)

		time.Sleep(time.Duration(cfg.Interval.Availability) * time.Second)
	}
}


func MonitorPorts(cfg *config.Config, serverIp string, port int, wg *sync.WaitGroup){
	defer wg.Done()
	for {
		status := 0
		address := fmt.Sprintf("%s:%d", serverIp, port)
		conn, err := net.DialTimeout("tcp", address, 5*time.Second)
		if err != nil {
			logger.Error("Failed to connect to port:", port, "on IP:", serverIp, err)
		} else {
			status = 1
			conn.Close()
		}
		logger.Info("Port Availability for ip:", serverIp, "port:", port, status)

		time.Sleep(time.Duration(cfg.Interval.Port) * time.Second)
	}
}


func MonitorService(cfg *config.Config, serviceName string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		cmd := exec.Command("systemctl", "is-active", "--quiet", serviceName)
		err := cmd.Run()
		if err != nil {
			logger.Info("Service is not running:", serviceName)
			restartCmd := exec.Command("systemctl", "restart", serviceName)
			err := restartCmd.Run()
			if err != nil {
				logger.Error("Failed to restart:", serviceName, err)
			} else {
				logger.Info("Service restarted successfully:", serviceName)
			}
		} else {
			logger.Info("Service is running:", serviceName)
		}

		time.Sleep(time.Duration(cfg.Interval.Service) * time.Second)
	}
}


func StartServerMonitoring(cfg *config.Config, wg *sync.WaitGroup){
	for _, server := range cfg.Monitoring.Servers {
		wg.Add(1)
		go MonitorAvailability(cfg, server, wg)
		for _, port := range server.Ports {
			wg.Add(1)
			go MonitorPorts(cfg, server.ServerIp, port, wg)
		}
	}
}


func StartSelfMonitoring(cfg *config.Config, wg *sync.WaitGroup){
	localhost := config.Server{
		ServerIp: "localhost",
		Ports: cfg.SelfMonitoring.Ports,
		Services: cfg.SelfMonitoring.Services,
	}

	wg.Add(1)
	go MonitorAvailability(cfg, localhost, wg)
	
	for _, port := range localhost.Ports {
		wg.Add(1)
		go MonitorPorts(cfg, localhost.ServerIp, port, wg)
	}

	for _, service := range localhost.Services {
		wg.Add(1)
		go MonitorService(cfg, service, wg)
	}
}