package file_service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var allFiles []string
var allReadFiles [][]byte

func ReadSomeFiles() [][]byte {

	maxLen := 0
	minLen := 100000000

	getDirFiles("/Users/wsl/Projects/go/src/github.com/docker")
	getDirFiles("/Users/wsl/Projects/go/src/github.com/aliyun")

	for _, name := range allFiles {
		content, err := read(name)
		if err != nil {
			fmt.Printf("%+v\n", err)
		} else {
			ll := len(content)
			if ll > maxLen {
				maxLen = ll
			}
			if ll < minLen {
				minLen = ll
			}
			fmt.Println(len(content), name)
			allReadFiles = append(allReadFiles, content)
		}
	}
	fmt.Println(maxLen, minLen)
	return allReadFiles
}

func read(filePath string) ([]byte, error) {
	fi, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	return ioutil.ReadAll(fi)
}

func getDirFiles(rPath string) {
	files, err := ioutil.ReadDir(rPath)
	if err != nil {
		fmt.Printf("error read %s : %+v", rPath, err)
		return
	}
	for _, file := range files {
		allPath := path.Join(rPath, file.Name())
		if file.IsDir() {
			getDirFiles(allPath)
		} else {
			if len(allFiles) == 10000 {
				return
			}
			if strings.HasSuffix(file.Name(), ".go") {
				allFiles = append(allFiles, allPath)
			}
		}
	}
}
