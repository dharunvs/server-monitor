package connection

import (
	"database/sql"
	"fmt"
	"sync"

	"root/config"
	"root/logger"
)

type MonitoringData struct {
	Host		string
	Type		string
	Parameter	string
	Value		string
}

var MonitoringDataChannel chan MonitoringData

func init () {
	MonitoringDataChannel = make(chan MonitoringData)
}

func WriteToDB(db *sql.DB, dataChannel <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for msg := range dataChannel {
		logger.Info(fmt.Sprintf("TODO: %v", msg))
	}
}

func WriteToMonitoringData(db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()

	query := "INSERT INTO monitoring_data (host, type, parameter, value) VALUES ($1, $2, $3, $4)"

	for data := range MonitoringDataChannel {
		_, err := db.Exec(query, data.Host, data.Type, data.Parameter, data.Value)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in WriteToMonitoringData(): %v", err))
		}
	}
}


func GetDatabase(database *config.Database) (*sql.DB, error){
	connString := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s", database.Host, database.Port, database.User, database.Password, database.Database)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}