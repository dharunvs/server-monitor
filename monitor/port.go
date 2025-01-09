package monitor

import (
	"fmt"
	"net"
	"sync"
	"time"

	"root/config"
	"root/logger"
)

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
