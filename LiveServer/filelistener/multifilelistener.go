package filelistener

func NewMultiFileListener() MultiFileListener {
	mfl := MultiFileListener{}
	mfl.listenChan = make(chan struct{})

	return mfl
}

type MultiFileListener struct {
	filePaths  []string
	listenChan chan struct{}
}

func (mfl *MultiFileListener) Add(fileListener FileListener) {
	if !mfl.alreadyAdded(&fileListener) {
		mfl.filePaths = append(mfl.filePaths, fileListener.filePath)
		fileListener.listenChan = mfl.listenChan
		fileListener.Start()
	}
}

func (mfl *MultiFileListener) GetListenChan() <-chan struct{} {
	return mfl.listenChan
}

func (mfl *MultiFileListener) alreadyAdded(fileListener *FileListener) bool {
	filePath := fileListener.filePath

	for _, fp := range mfl.filePaths {
		if filePath == fp {
			return true
		}
	}

	return false
}
