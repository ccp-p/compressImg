package auto

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"test/logger"

	"github.com/fsnotify/fsnotify"
)

//imgObj := map[string]string{"url": location, "fileName": "1.png"}

func Downloaded(urls []map[string]string) {

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, imgObj := range urls {

		go func(imgObj map[string]string) {
			println(imgObj)
			logger.Logger.Println(imgObj)
			defer wg.Done()
			downloadPng(imgObj)
		}(imgObj)
	}

	wg.Wait()
	fmt.Println("All images downloaded and saved successfully.")
	logger.Logger.Println("All images downloaded and saved successfully.")
}

func downloadPng(imgObj map[string]string) {
	fmt.Println("Downloading png from", imgObj["url"])
	logger.Logger.Println("Downloading png from", imgObj["url"])
	response, err := http.Get(imgObj["url"])
	if err != nil {
		fmt.Println("Error downloading png:", err)
		return
	}
	defer response.Body.Close()

	absPath := imgObj["absPath"]

	fileName := filepath.Base(absPath)
	// 如果源图片名包含@2x，去掉这个后缀
	fileName = strings.Replace(fileName, "@2x", "", -1)

	dirPath := filepath.Dir(absPath)

	newPath := filepath.Join(filepath.Join(dirPath, "\\compressed"))

	// 检查目录是否存在
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		// 目录不存在，创建它
		err := os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
		}
	}

	file, err := os.Create(filepath.Join(newPath, fileName))
	if err != nil {
		fmt.Println("Error creating file:", err)
		logger.Logger.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("Error saving png:", err)
		return
	}

	log.Println("Png downloaded and saved successfully:", fileName)
}

func getFileName(url string) string {
	// Extract the file name from the URL
	// For example, "https://shadow.elemecdn.com/app/element/hamburger.9cf7b091-55e9-11e9-a976-7f4d0b07eef6.png" will return "hamburger.png"
	fileName := url[strings.LastIndex(url, "/")+1:]

	return fileName
}

func IsImage(filePath string) bool {
	ext := filepath.Ext(filePath)
	switch ext {
	case ".png", ".jpg", ".jpeg":
		return true
	}
	return false

}
func WatchFolder(dirPath string) (<-chan string, error) {
	watcher, err := fsnotify.NewWatcher()
	fmt.Println("watcher:", watcher)
	logger.Logger.Println("watcher:", watcher)
	if err != nil {
		return nil, err
	}

	files := make(chan string)
	go func() {
		defer close(files)
		for {
			select {
			case event, ok := <-watcher.Events:
				fmt.Println("event====", event)
				logger.Logger.Println("event===", event)
				if !ok {
					return
				}
				//创建文件
				stand := event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write
				if stand {
					files <- event.Name
					fmt.Println("Modified file:", event.Name)
					logger.Logger.Println("Modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				fmt.Println("error:", err)
				logger.Logger.Println("error:", err)
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(dirPath)
	if err != nil {
		return nil, err
	}

	return files, nil
}
