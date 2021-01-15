package dir

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/hiromaily/go-gin-wrapper/pkg/regexps"
)

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
		if regexps.IsInvisiblefile(fi.Name()) {
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
				if regexps.IsExtFile(fi.Name(), ex) {
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
