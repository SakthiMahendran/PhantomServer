package filelistener

//First read the "filelistener.go", description and then the code line by line with the comments to for a better understanding.

//Description
/*
	"multifilelistener" is designed to handle multiple "filelistener".
	it sets a single channel for all of the "filelistener" in it's control to send signal.
	So that we get a single channel to look for signal instead of looking
	into all of the filelistener's channel.
*/

//Makes a new "MultiFileListener".
func NewMultiFileListener() MultiFileListener {
	mfl := MultiFileListener{}           //Instantiation.
	mfl.listenChan = make(chan struct{}) //Makes a new channel for signaling changes in any file.

	return mfl // Returning.
}

type MultiFileListener struct {
	fileListeners []*FileListener // Contains the paths of files to be listened for changes.
	listenChan    chan struct{}   // Singnal will be send through this channel if there is changes in any file.
}

//Adds a new "FileListener" to be handled.
func (mfl *MultiFileListener) Add(fileListener *FileListener) {
	if !mfl.alreadyAdded(fileListener) { // Checks wheather the file is already added.
		//if not
		mfl.fileListeners = append(mfl.fileListeners, fileListener) //Add the path of the file to filePaths.
		fileListener.listenChan = mfl.listenChan                    //Sets the signaling channel for filelistener.
		fileListener.Start()                                        //Starts the file listener.
	}
}

func (mfl *MultiFileListener) GetListenChan() <-chan struct{} {
	return mfl.listenChan // returns the listenChan
}

//Checks wheather the given fileListener is already added.
func (mfl *MultiFileListener) alreadyAdded(fileListener *FileListener) bool {
	filePath := fileListener.filePath // Gets the "filePath" of fileListener.

	for _, fp := range mfl.fileListeners { //Checks wheather the filePath is in "filePaths" slice.
		if filePath == fp.filePath {
			return true // return "true" if yes.
		}
	}

	return false // else return "false".
}

func (mfl *MultiFileListener) Reset() {
	for _, fl := range mfl.fileListeners {
		fl.Stop()
	}

	mfl.fileListeners = nil
	mfl.fileListeners = make([]*FileListener, 0)
}
