package modules

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

var pages map[string][]byte

func getPages() (map[string][]byte, error) {
	if pages != nil {
		return pages, nil
	}
	pages := make(map[string][]byte)
	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	pagesDir := path.Join(root, "./pages")
	files, err := ioutil.ReadDir(pagesDir)
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}
		filename := fileInfo.Name()
		filePath := path.Join(pagesDir, filename)
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		pages[filename] = fileContent
	}
	return pages, nil
}

func GetPageByStatus(status int) ([]byte, error) {
	pages, err := getPages()
	if err != nil {
		return nil, err
	}
	return pages[strconv.Itoa(status)+".html"], nil

}
