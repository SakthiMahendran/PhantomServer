package filelistener

import (
	"os"
	"time"
)

// RATE defines the rate at which the file is checked for changes (i.e Now the file is checked on every 250ms for changes).
const RATE time.Duration = time.Millisecond * 250

// FileListener struct defines the properties of the file listener
type FileListener struct {
	filePath   string          //filePath defines the path of the file to be listened for changes.
	rate       time.Duration   //rate defines the rate of listening for changes.
	listenChan chan<- struct{} //if there is any changes in the file then signal will be sent through this channel.
	stoped     bool            //stoped is a boolean flag that indicates whether the file listener has been stopped or not.
}

// NewFileListener returns a new FileListener object
func NewFileListener(filePath string) FileListener {
	// Initialize a new FileListener object with default values
	fl := FileListener{
		filePath: filePath, //Set the path of the file to be listened
		rate:     RATE,     //Set the rate of listening to file changes
		stoped:   false,    //Set the initial value of stopped flag
	}
	//Note: the "fl.listenChan" will be setted by "MultiFileListener".

	return fl // Return the FileListener Object.
}

// Start starts the file listener
func (fl *FileListener) Start() {
	//Starts a new goroutine which will listen for the changes in the file.
	go func() {
		initStat, _ := os.Stat(fl.filePath) // Get the initial state of the file
		initMod := initStat.ModTime()       // Get the initial "Last Modified Date" of the file.

		ticker := time.NewTicker(fl.rate) // Start a new ticker based on the file listener rate

		for range ticker.C { // Wait for the ticker to fire
			if fl.stoped { // Check if the file listener has been stopped
				ticker.Stop() // Stop the ticker
				return        // Return from the goroutine
			}

			lastStat, _ := os.Stat(fl.filePath) // Get the current state of the file
			lastMod := lastStat.ModTime()       // Get the current "Last Modified Date" of the file.

			if (lastMod != initMod) && !(fl.stoped) { // Check whether the file is modified.
				fl.listenChan <- struct{}{} // If the file has been modified, send a signal through the channel
				initMod = lastMod           // Set current "Last Modified Date" as initial "Last Modified Date"
			}
		}
	}()
}

// Stop stops the file listener
func (fl *FileListener) Stop() {
	fl.stoped = true // Set the stoped flag to true to stop the file listener
}
