package auto

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

//imgObj := map[string]string{"url": location, "fileName": "1.png"}

func Downloaded(urls []map[string]string) {

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, imgObj := range urls {

		go func(imgObj map[string]string) {
			println(imgObj)
			defer wg.Done()
			downloadPng(imgObj)
		}(imgObj)
	}

	wg.Wait()
	fmt.Println("All images downloaded and saved successfully.")
}

func downloadPng(imgObj map[string]string) {
	fmt.Println("Downloading png from", imgObj["url"])
	response, err := http.Get(imgObj["url"])
	if err != nil {
		fmt.Println("Error downloading png:", err)
		return
	}
	defer response.Body.Close()

	fileName := getFileName(imgObj["fileName"])
	newPath := "../png/"
	// 检查目录是否存在
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		// 目录不存在，创建它
		err := os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
		}
	}
	file, err := os.Create(newPath + fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("Error saving png:", err)
		return
	}

	fmt.Println("Png downloaded and saved successfully:", fileName)
}

func getFileName(url string) string {
	// Extract the file name from the URL
	// For example, "https://shadow.elemecdn.com/app/element/hamburger.9cf7b091-55e9-11e9-a976-7f4d0b07eef6.png" will return "hamburger.png"
	fileName := url[strings.LastIndex(url, "/")+1:]

	return fileName
}
