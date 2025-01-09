package monitor

import (
	"os/exec"
	"sync"
	"time"

	"root/config"
	"root/connection"
	"root/logger"
)


func MonitorService(cfg *config.Config, serviceName string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		var status string;
		cmd := exec.Command("systemctl", "is-active", "--quiet", serviceName)
		err := cmd.Run()
		if err != nil {
			logger.Info("Service is not running:", serviceName)
			restartCmd := exec.Command("systemctl", "restart", serviceName)
			err := restartCmd.Run()
			if err != nil {
				logger.Error("Failed to restart:", serviceName, err)
				status = "Not Running"

			} else {
				logger.Info("Service restarted successfully:", serviceName)
				status = "Restarted"
			}
		} else {
			logger.Info("Service is running:", serviceName)
			status = "Running"
		}
		
		data := connection.MonitoringData{
			Host: "localhost",
			Type: "Service",
			Parameter: serviceName,
			Value: status,
		}

		connection.MonitoringDataChannel <- data;
		time.Sleep(time.Duration(cfg.Interval.Service) * time.Second)
	}
}