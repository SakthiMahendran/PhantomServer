package filelistener

import (
	"fmt"
	"os"
	"time"
)

func NewFileListener(filePath string, rate time.Duration) FileListener {
	fl := FileListener{}

	fl.filePath = filePath
	fl.rate = rate

	return fl
}

type FileListener struct {
	filePath   string          //filePath defines the path of the file to be listened for changes.
	rate       time.Duration   //rate defines the rate of listening for changes (if rate is 250ms then file checked for changes on every 250ms).
	listenChan chan<- struct{} //if there is any changes in file then signal is send on this channel.
}

//Starts the file listening
func (fl *FileListener) Start() {
	go func(lstnChan chan<- struct{}, filePath string, rate time.Duration) {

		initStat, _ := os.Stat(filePath)
		initMod := initStat.ModTime()

		for range time.Tick(rate) {

			lastStat, err := os.Stat(filePath)

			if err != nil {
				fmt.Printf(err.Error())
			}

			lastMod := lastStat.ModTime()

			if lastMod != initMod {
				lstnChan <- struct{}{}
				initMod = lastMod
			}
		}
	}(fl.listenChan, fl.filePath, fl.rate)
}
