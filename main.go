package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

func main() {
	if len(os.Args) <= 2 {
		fmt.Println("please enter root directory path and string to search")
	}
	if len(os.Args[2]) < 2 {
		fmt.Println("minimum search string should contain at least 2 characters")
	}
	if exists, err := isDirectoryExists(os.Args[1]); err != nil {
		fmt.Println(err)
	} else if !exists {
		fmt.Println("directory does not exists")
	} else {
		searchRecursively(os.Args[1], (os.Args[2]))
	}
}

func searchRecursively(path string, s string) {
	files := getAllFiles(path)
	if files != nil {
		for path, ext := range getAllFiles(path) {
			fmt.Printf("searching %s...\n", path)
			switch ext {
			case ".pdf":
				if found, _ := searchPDF(path, s); found {
					fmt.Printf("FOUND! %s\nFOUND! %s\nFOUND! %s\n", path, path, path)
				}
				break

			case ".docx":
				if found, _ := searchDocx(path, s); found {
					fmt.Printf("FOUND! %s\nFOUND! %s\nFOUND! %s\n", path, path, path)
				}
				break

			case ".doc":
				break
			}
		}
	}
}

func getAllFiles(path string) map[string]string {
	fileList := make(map[string]string, 0)
	if err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if strings.EqualFold(ext, ".pdf") || strings.EqualFold(ext, ".doc") || strings.EqualFold(ext, ".docx") {
			fileList[path] = strings.ToLower(ext)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
		return nil
	}
	return fileList
}

func isDirectoryExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err

}

func printError() {
	err := recover()
	if err != nil {
		fmt.Println(err)
	}
}

func searchPDF(path string, s string) (bool, error) {
	defer printError()
	_, r, err := pdf.Open(path)
	if err != nil {
		return false, err
	}
	totalPage := r.NumPage()
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		content := p.Content()
		var textBuilder bytes.Buffer
		defer textBuilder.Reset()
		if content.Text != nil {
			for _, t := range content.Text {
				textBuilder.WriteString(t.S)
			}
		}
		str := textBuilder.String()
		if strings.Contains(str, s) {
			return true, nil
		}

	}
	return false, nil
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
func searchDocx(path string, s string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	body, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		fmt.Println(err)
	}
	for _, zipFile := range zipReader.File {
		if zipFile.Name == "word/document.xml" {
			unzippedFileBytes, err := readZipFile(zipFile)

			if err != nil {
				log.Println(err)

			}
			if bytes.Contains(unzippedFileBytes, []byte(s)) {
				return true, nil
			}
			continue
		}
	}

	// for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
	// 	p := r.Page(pageIndex)
	// 	if p.V.IsNull() {
	// 		continue
	// 	}
	// 	content := p.Content()
	// 	var textBuilder bytes.Buffer
	// 	defer textBuilder.Reset()
	// 	if content.Text != nil {
	// 		for _, t := range content.Text {
	// 			textBuilder.WriteString(t.S)
	// 		}
	// 	}
	// 	str := textBuilder.String()
	// 	if strings.Contains(str, s) {
	// 		return true, nil
	// 	}

	// }
	return false, nil
}
