package filelistener

import (
	"os"
	"time"
)

//First read the description and then the code line by line with the comments to for a better understanding.

//Description
/*
"filelistener" will listen at a particular file at a particular rate for changes and
then send a signal in a channel (listenChan) if there is any changes made to the file.
*/

//RATE defines the rate at which the file is checked for changes (i.e Now the file is checked on every 250ms for changes).
const RATE time.Duration = time.Millisecond * 250

//Makes a new FileListener.
func NewFileListener(filePath string) FileListener {
	fl := FileListener{} //Instantiation.

	fl.filePath = filePath //Setting the filePath (Path of the file to be listened).
	fl.rate = RATE         //Setting the rate.
	//Note: the "fl.listenChan" will be setted by "MultiFileListener".

	return fl //Returning the FileListener Object.
}

type FileListener struct {
	filePath   string          //filePath defines the path of the file to be listened for changes.
	rate       time.Duration   //rate defines the rate of listening for changes.
	listenChan chan<- struct{} //if there is any changes in the file then signal will be send through this channel.
}

//Starts the file listening
func (fl *FileListener) Start() {
	//Starts a new goroutine which will listen for the changes in the file.
	go func(lstnChan chan<- struct{}, filePath string, rate time.Duration) {

		initStat, _ := os.Stat(filePath) //Getting the initial file state.
		initMod := initStat.ModTime()    //Getting the initial "Last Modified Date" of the file.

		for range time.Tick(rate) { //Waiting for a time duration that is given by "rate" (Now it waits for 250ms).

			lastStat, _ := os.Stat(filePath) //Getting the current file state.
			lastMod := lastStat.ModTime()    //Getting the current "Last Modified Date" of the file.

			if lastMod != initMod { //Cheacking wheather the file is modified.
				lstnChan <- struct{}{} //If modified sending singnal through the channel
				initMod = lastMod      //Setting current "Last Modified Date" as initial "Last modified date".
			}
		}
	}(fl.listenChan, fl.filePath, fl.rate)
}
