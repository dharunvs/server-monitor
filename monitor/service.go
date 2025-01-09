package monitor

import (
	"os/exec"
	"sync"
	"time"

	"root/config"
	"root/logger"
)


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