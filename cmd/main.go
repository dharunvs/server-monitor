package main

import (
	"os"
	"sync"

	"root/config"
	"root/connection"
	"root/logger"
	"root/monitor"
	"root/notifier"

	_ "github.com/lib/pq"
	// "root/backup"
)

func main() {
	logger.Info("Starting MonitoringModule")
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		logger.Error("Error loading config:", err)
		return
	}

	notificationChannel := make(chan string)

	db, err := connection.GetDatabase(&cfg.Database)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	defer db.Close()

	var wg sync.WaitGroup

	wg.Add(3)
	go notifier.StartNotifier(notificationChannel, &wg)

	tableChannels := make(map[string]chan string)
	tableChannels["monitoringdata"] = make(chan string)

	// connection.WriteToDB()
	go monitor.StartServerMonitoring(cfg, &wg)
	go monitor.StartSelfMonitoring(cfg, &wg)

	// file, err := backup.DumpDatabase(&cfg.Backup.SourceDB, "monitoringmodule", "/home/dharunvs/Documents/serverMonitoring/bkps/")
	// err = backup.DumpAndRestore(&cfg.Backup.SourceDB, &cfg.Backup.DestinationDBs[0], "vs_test_database","/home/dharunvs/Documents/serverMonitoring/bkps/")
	// if err != nil {
	// 	logger.Error(err)
	// }

	wg.Wait()

	// close(tableChannels["monitoringdata"] )
	// logger.Info("MonitoringModule ended", connString)
}
