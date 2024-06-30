package connection

import (
	"fmt"
	"sync"
	"database/sql"

	"root/logger"
	"root/config"
)

func WriteToDB(db *sql.DB, table string, dataChannel <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for msg := range dataChannel {
		_, err := db.Exec("INSERT INTO "+table+" (value) VALUES ($1)", msg)
		if err != nil {
			logger.Error("Failed to insert message into "+table+":", err)
		} else {
			logger.Info("Message inserted into "+table+":", msg)
		}
	}
}

func SelectFromDB() {
	
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