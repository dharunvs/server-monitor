package monitor

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-ping/ping"

	"root/config"
	"root/connection"
	"root/logger"
	"root/notifier"
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

		data := connection.MonitoringData{
			Host: server.ServerIp,
			Type: "Server",
			Parameter: "Availability",
			Value: strconv.Itoa(availability),
		}
		logger.Info("Availability for ip:", server.ServerIp, availability)
		connection.MonitoringDataChannel <- data;

		if availability == 0 {
			notifier.NotificationDataChannel <- notifier.CreateNotification(fmt.Sprintf("Server not available: %v", server.ServerIp))
		}

		time.Sleep(time.Duration(cfg.Interval.Availability) * time.Second)
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