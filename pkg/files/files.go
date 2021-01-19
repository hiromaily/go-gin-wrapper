package files

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sync"
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
	filelist, err := ioutil.ReadDir(basePath)
	if err != nil {
		// fmt.Printf("error : %s\n", err)
		return
	}

	for _, file := range filelist {
		if IsInvisiblefile(file.Name()) {
			continue
		}

		fullPath := filepath.Join(basePath, file.Name())
		if file.IsDir() {
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
			for _, ext := range exts {
				if IsExtFile(file.Name(), ext) {
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

// IsStaticFile checks whether target is static file which has extension in file name or not
func IsStaticFile(target string) bool {
	//is there any suffix
	//.+\.(csv|pdf)
	return compileRegexp(`^.*\.`, target)
}

// IsInvisiblefile checks whether target is invisible file or not
func IsInvisiblefile(target string) bool {
	return compileRegexp(`^[\\.].*$`, target)
}

// IsExtFile checks whether target includes ext string
func IsExtFile(target, ext string) bool {
	return compileRegexp(`^.*\.`+ext+`$`, target)
}

// compileRegexp compiles regex string
func compileRegexp(reg, str string) bool {
	return regexp.MustCompile(reg).Match([]byte(str))
}
