package files

import (
	reg "github.com/hiromaily/golibs/regexp"
	"io/ioutil"
	"path/filepath"
	"sync"
)

var fileNames []string

// checkDirectory is to check directory using goroutine as semaphore
func checkDirectory(basePath string, extensions []string, ch chan<- string, chSmp chan bool, wg *sync.WaitGroup, closeFlg bool) {

	// read directory
	fis, err := ioutil.ReadDir(basePath)
	if err != nil {
		//fmt.Printf("error : %s\n", err)
		return
	}

	for _, fi := range fis {
		//fmt.Printf("file name is %s\n", fi.Name())
		if reg.IsInvisiblefile(fi.Name()) {
			continue
		}

		fullPath := filepath.Join(basePath, fi.Name())
		//fmt.Printf("full path is %s\n", fullPath)

		if fi.IsDir() {
			wg.Add(1)
			chSmp <- true

			//fmt.Println("this is directory. skip.")
			//check more deep directory
			go func() {
				defer func() {
					<-chSmp
					wg.Done()
				}()
				checkDirectory(fullPath, extensions, ch, chSmp, wg, false)
			}()
		} else {
			for _, ex := range extensions {
				//fmt.Printf("search %s from %s\n", ex, fi.Name())
				if reg.IsExtFile(fi.Name(), ex) {
					ch <- fullPath
				}
			}
		}
	}

	if closeFlg {
		wg.Wait()
		close(ch)
	}
}

// checkDirectorySingle is check directory using single goroutine
func checkDirectorySingle(basePath string, extensions []string, ch chan<- string, closeFlg bool) {
	// read directory
	fis, err := ioutil.ReadDir(basePath)
	if err != nil {
		//fmt.Printf("error : %s\n", err)
		return
	}

	for _, fi := range fis {
		//fmt.Printf("file name is %s\n", fi.Name())
		if reg.IsInvisiblefile(fi.Name()) {
			continue
		}

		fullPath := filepath.Join(basePath, fi.Name())
		//fmt.Printf("full path is %s\n", fullPath)

		if fi.IsDir() {
			//fmt.Println("this is directory. skip.")

			//check more deep directory
			checkDirectorySingle(fullPath, extensions, ch, false)
		} else {
			for _, ex := range extensions {
				//fmt.Printf("search %s from %s\n", ex, fi.Name())
				if reg.IsExtFile(fi.Name(), ex) {
					ch <- fullPath
				}
			}
		}
	}
	if closeFlg {
		close(ch)
	}
}

// checkDirectoryJIC is check directory without goroutine
func checkDirectoryJIC(basePath string, extensions []string) {
	// read directory
	fis, err := ioutil.ReadDir(basePath)
	if err != nil {
		//fmt.Printf("error : %s\n", err)
		return
	}

	for _, fi := range fis {
		//fmt.Printf("file name is %s\n", fi.Name())
		if reg.IsInvisiblefile(fi.Name()) {
			continue
		}

		fullPath := filepath.Join(basePath, fi.Name())
		//fmt.Printf("full path is %s\n", fullPath)

		if fi.IsDir() {
			//fmt.Println("this is directory. skip.")

			//check more deep directory
			checkDirectoryJIC(fullPath, extensions)
		} else {
			for _, ex := range extensions {
				//fmt.Printf("search %s from %s\n", ex, fi.Name())
				if reg.IsExtFile(fi.Name(), ex) {
					fileNames = append(fileNames, fullPath)
				}
			}
		}
	}
}

// GetFileList is to get file list using multiple goroutine
func GetFileList(basePath string, extensions []string) []string {
	files := []string{}
	pool := 20
	ch := make(chan string, pool)
	chSmp := make(chan bool, pool)
	wg := &sync.WaitGroup{}

	go checkDirectory(basePath, extensions, ch, chSmp, wg, true)
	for {
		v, ok := <-ch
		if ok {
			files = append(files, v)
		} else {
			break
		}
	}
	return files
}

// GetFileListSingle is to get file list using a goroutine
func GetFileListSingle(basePath string, extensions []string) []string {
	files := []string{}
	ch := make(chan string)
	go checkDirectorySingle(basePath, extensions, ch, true)
	for {
		v, ok := <-ch
		if ok {
			files = append(files, v)
		} else {
			break
		}
	}
	return files
}

// GetFileListJIC is to get file list without goroutine
// It doesn't use channel. But faster. (JIC means Just In Case)
func GetFileListJIC(basePath string, extensions []string) []string {
	fileNames = []string{}

	checkDirectoryJIC(basePath, extensions)
	return fileNames
}
