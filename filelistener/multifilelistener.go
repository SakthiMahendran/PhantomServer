package filelistener

// MultiFileListener handles multiple FileListeners by setting a single channel
// for all of them to send signal through. This allows for a single channel to be
// used instead of looking into each individual FileListener's channel.
type MultiFileListener struct {
	fileListeners []*FileListener // contains the paths of files to be listened for changes
	listenChan    chan struct{}   // signal will be sent through this channel if there are changes in any file
}

// NewMultiFileListener creates a new MultiFileListener instance and returns it.
func NewMultiFileListener() MultiFileListener {
	mfl := MultiFileListener{}
	mfl.listenChan = make(chan struct{})
	return mfl
}

// Add adds a new FileListener to be handled by the MultiFileListener.
func (mfl *MultiFileListener) Add(fileListener *FileListener) {
	if !mfl.alreadyAdded(fileListener) {
		mfl.fileListeners = append(mfl.fileListeners, fileListener)
		fileListener.listenChan = mfl.listenChan
		fileListener.Start()
	}
}

// GetListenChan returns the listen channel.
func (mfl *MultiFileListener) GetListenChan() <-chan struct{} {
	return mfl.listenChan
}

// alreadyAdded checks whether the given fileListener is already added.
func (mfl *MultiFileListener) alreadyAdded(fileListener *FileListener) bool {
	filePath := fileListener.filePath
	for _, fp := range mfl.fileListeners {
		if filePath == fp.filePath {
			return true
		}
	}
	return false
}

// Reset stops all the file listeners and clears the fileListeners slice.
func (mfl *MultiFileListener) Reset() {
	for _, fl := range mfl.fileListeners {
		fl.Stop()
	}
	mfl.fileListeners = nil
	mfl.fileListeners = make([]*FileListener, 0)
}
