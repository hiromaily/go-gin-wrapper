package dir

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/hiromaily/go-gin-wrapper/pkg/regexps"
)

// GetFileList returns file path from basePath directory with acceptable extension files
func GetFileList(basePath string, exts []string) []string {
	files := []string{}
	pool := 20
	ch := make(chan string, pool)
	chSmp := make(chan bool, pool)
	wg := &sync.WaitGroup{}

	go checkDir(basePath, exts, ch, chSmp, wg, true)
	for {
		v, ok := <-ch
		if !ok {
			break
		}
		files = append(files, v)
	}
	return files
}

func checkDir(basePath string, exts []string, ch chan<- string, chSmp chan bool, wg *sync.WaitGroup, isClose bool) {
	fis, err := ioutil.ReadDir(basePath)
	if err != nil {
		// fmt.Printf("error : %s\n", err)
		return
	}

	for _, fi := range fis {
		if regexps.IsInvisiblefile(fi.Name()) {
			continue
		}

		fullPath := filepath.Join(basePath, fi.Name())
		if fi.IsDir() {
			wg.Add(1)
			chSmp <- true

			// check more deep directory
			go func() {
				defer func() {
					<-chSmp
					wg.Done()
				}()
				checkDir(fullPath, exts, ch, chSmp, wg, false)
			}()
		} else {
			for _, ex := range exts {
				if regexps.IsExtFile(fi.Name(), ex) {
					ch <- fullPath
				}
			}
		}
	}

	if isClose {
		wg.Wait()
		close(ch)
	}
}
