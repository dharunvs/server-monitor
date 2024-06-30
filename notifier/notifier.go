package notifier

import (
	"sync"
	"root/logger"
)

type Notification struct {
	Message		string		`json:"message"`
}

func StartNotifier(notifications <- chan string, wg *sync.WaitGroup){
	defer wg.Done()
	for notfication := range notifications {
		SendNotification(notfication)
	}
}

func SendNotification(notfication string){
	logger.Info("TODO: send notificaiton:", notfication)
}