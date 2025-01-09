package main

import (
	"os"
	"sync"

	_ "github.com/lib/pq"

	"root/config"
	"root/connection"
	"root/logger"
	"root/monitor"
	"root/notifier"
	// "root/backup"
)

var wg sync.WaitGroup


func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		logger.Error("Error in config.LoadConfig():", err)
		return
	}

	wg.Add(1)
	go notifier.StartNotifier(&cfg.Notifier.Telegram, &wg)

	
	db, err := connection.GetDatabase(&cfg.Database)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	logger.Info("Connected to database", cfg.Database)
	defer db.Close()


	logger.Info("Starting Server Monitor")
	notifier.NotificationDataChannel <- notifier.CreateNotification("Starting Server Monitor")

	wg.Add(1)
	go connection.WriteToMonitoringData(db, &wg)

	wg.Add(2)
	go monitor.StartServerMonitoring(cfg, &wg)
	go monitor.StartSelfMonitoring(cfg, &wg)

	// file, err := backup.DumpDatabase(&cfg.Backup.SourceDB, "monitoringmodule", "/home/dharunvs/Documents/serverMonitoring/bkps/")
	// err = backup.DumpAndRestore(&cfg.Backup.SourceDB, &cfg.Backup.DestinationDBs[0], "vs_test_database","/home/dharunvs/Documents/serverMonitoring/bkps/")
	// if err != nil {
	// 	logger.Error(err)
	// }

	wg.Wait()

	logger.Info("Server Monitor ended")
}
